import json
import os
import sys
import requests
import hashlib
import base64
import bech32
import cosmospy

args = sys.argv[1:]

# Get the args
build_dir = args[0]
genesis_file = f"{build_dir}/node0/desmos/config/genesis.json"

chain_state_url = args[2]
chain_state_file = f"{build_dir}/state.json"
output_file = f"{build_dir}/output_state.json"

# Get the chain state inside the build dir
with requests.get(chain_state_url) as r, open(chain_state_file, 'w') as f:
    f.write(json.dumps(r.json()))

with open(chain_state_file, 'r') as chain_state_f, open(genesis_file, 'r') as genesis_f, open(output_file, 'w') as out:
    chain_state = json.load(chain_state_f)
    genesis = json.load(genesis_f)

    # -------------------------------
    # --- Update the bank state
    modules_addresses = [
        'desmos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8n8fv78',
        'desmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3prylw0',
        'desmos1tygms3xhhs3yv487phx3dw4a95jn7t7l4rcwcm',
    ]

    genesis['app_state']['bank']['supply'] = []
    for balance in chain_state['app_state']['bank']['balances']:
        if balance['address'] in modules_addresses:
            balance['coins'] = []

        genesis['app_state']['bank']['balances'].append(balance)

    # -------------------------------
    # --- Fix modules state

    # x/auth
    genesis['app_state']['auth']['accounts'] += chain_state['app_state']['auth']['accounts']

    # x/gov
    genesis['app_state']['gov']['deposit_params']['max_deposit_period'] = '120s'
    genesis['app_state']['gov']['voting_params']['voting_period'] = '120s'
    genesis['app_state']['gov']['deposit_params']['min_deposit'] = [{'amount': '10000000', 'denom': 'udaric'}]

    # -------------------------------
    # --- Copy modules state over
    genesis['app_state']['crisis'] = chain_state['app_state']['crisis']
    genesis['app_state']['ibc'] = chain_state['app_state']['ibc']
    genesis['app_state']['profiles'] = chain_state['app_state']['profiles']

    custom_modules = ['profiles', 'subspaces', 'relationships', 'wasm']
    for module in custom_modules:
        if module in chain_state['app_state']:
            genesis['app_state'][module] = chain_state['app_state'][module]
        else:
            del (genesis['app_state'][module])

    # -------------------------------
    # --- Write the file

    out.write(json.dumps(genesis))
    os.system(f"sed -i 's/\"stake\"/\"udaric\"/g' {output_file}")
    os.system(f"sed -i 's/\"udsm\"/\"udaric\"/g' {output_file}")

nodes_amount = args[1]
for i in range(0, int(nodes_amount)):
    genesis_path = f"{build_dir}/node{i}/desmos/config/genesis.json"
    with open(genesis_path, 'w') as file:
        os.system(f"cp {output_file} {genesis_path}")
