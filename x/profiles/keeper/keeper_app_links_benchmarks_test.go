package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"

	"github.com/desmos-labs/desmos/v3/app"
	"github.com/desmos-labs/desmos/v3/testutil"
	"github.com/desmos-labs/desmos/v3/x/profiles/keeper"
	"github.com/desmos-labs/desmos/v3/x/profiles/types"
)

func setupBenchTest() (authkeeper.AccountKeeper, keeper.Keeper, sdk.Context) {
	// Define the store keys
	keys := sdk.NewKVStoreKeys(types.StoreKey, authtypes.StoreKey, paramstypes.StoreKey)
	tKeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	storeKey := keys[types.StoreKey]

	// Create an in-memory db
	memDB := db.NewMemDB()
	ms := store.NewCommitMultiStore(memDB)
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}
	for _, tKey := range tKeys {
		ms.MountStoreWithDB(tKey, sdk.StoreTypeTransient, memDB)
	}
	for _, memKey := range memKeys {
		ms.MountStoreWithDB(memKey, sdk.StoreTypeMemory, nil)
	}

	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	ctx := sdk.NewContext(ms, tmproto.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	cdc, legacyAminoCdc := app.MakeCodecs()

	paramsKeeper := paramskeeper.NewKeeper(
		cdc, legacyAminoCdc, keys[paramstypes.StoreKey], tKeys[paramstypes.TStoreKey],
	)

	ak := authkeeper.NewAccountKeeper(
		cdc,
		keys[authtypes.StoreKey],
		paramsKeeper.Subspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		app.GetMaccPerms(),
	)

	k := keeper.NewKeeper(
		cdc,
		legacyAminoCdc,
		storeKey,
		paramsKeeper.Subspace(types.DefaultParamsSpace),
		ak,
		nil,
		nil,
		nil,
		nil,
	)

	return ak, k, ctx
}

func generateRandomAppLinks(r *rand.Rand, linkNum int) []types.ApplicationLink {
	accounts := simtypes.RandomAccounts(r, r.Intn(linkNum))
	var appLinks []types.ApplicationLink
	for _, account := range accounts {
		link := types.NewApplicationLink(
			account.Address.String(),
			types.NewData(simtypes.RandStringOfLength(r, 5), simtypes.RandStringOfLength(r, 6)),
			types.ApplicationLinkStateInitialized,
			types.NewOracleRequest(
				0,
				1,
				types.NewOracleRequestCallData(simtypes.RandStringOfLength(r, 5), simtypes.RandStringOfLength(r, 5)),
				simtypes.RandStringOfLength(r, 10),
			),
			nil,
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			time.Date(2022, 1, 1, 00, 00, 00, 000, time.UTC),
		)

		appLinks = append(appLinks, link)
	}

	return appLinks
}

func BenchmarkKeeper_DeleteExpiredApplicationLinks(b *testing.B) {
	r := rand.New(rand.NewSource(100))
	ak, k, ctx := setupBenchTest()
	links := generateRandomAppLinks(r, 1)
	ctx, _ = ctx.CacheContext()

	for _, link := range links {
		ak.SetAccount(ctx, testutil.ProfileFromAddr(link.User))
		err := k.SaveApplicationLink(ctx, link)
		require.NoError(b, err)
	}

	b.Run("iterate and delete expired links", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.ReportAllocs()
			k.DeleteExpiredApplicationLinks(ctx)
		}
	})
}
