# Cosmos Blockchain Simulator

The Cosmos SDK offers a full fledged simulation framework to fuzz test every message defined by a module.

This functionality is provided by the[`SimApp`](https://github.com/hyperspeednetwork/hsnhub/blob/master/simapp/app.go),
which is a dummy application that is used for running the [`simulation`](https://github.com/hyperspeednetwork/hsnhub/tree/master/x/simulation) module.
This module defines all the simulation logic as well as the operations for randomized parameters like accounts, balances etc.

## Goals

The blockchain simulator tests how the blockchain application would behave under real life circumstances by generating and sending randomized messages.
The goal of this is to detect and debug failures that could halt a live chain, by providing logs and statistics about the operations run by the simulator as well as exporting the latest application state when a failure was found.

Its main difference with integration testing is that the simulator app allows you to pass parameters to customize the chain that's being simulated.
This comes in handy when trying to reproduce bugs that were generated in the provided operations (randomized or not).

## Simulation commands

The simulation app has different commands, each of which tests a different failure type:

- `AppImportExport`: The simulator exports the initial app state and then it creates a new app with the exported `genesis.json` as an input, checking for inconsistencies between the stores.
- `AppSimulationAfterImport`: Queues two simulations together. The first one provides the app state (_i.e_ genesis) to the second. Useful to test software upgrades or hard-forks from a live chain.
- `AppStateDeterminism`: Checks that all the nodes return the same values, in the same order.
- `BenchmarkInvariants`: Analyses the performance of running all the modules' invariants (_i.e_ secuentially runs a [benchmark](https://golang.org/pkg/testing/#hdr-Benchmarks) test). An invariant checks for differences between the values that are on the store and the passive tracker. Eg: total coins held by accounts vs total supply tracker.
- `FullAppSimulation`: General simulation mode. Runs the chain and the specified operations for a given number of blocks. Tests that there're no `panics` on the simulation. It does also run invariant checks on every `Period` but they are not benchmarked.

Each simulation must receive a set of inputs (_i.e_ flags) such as the number of blocks that the simulation is run, seed, block size, etc.
Check the full list of flags [here](https://github.com/hyperspeednetwork/hsnhub/blob/adf6ddd4a807c8363e33083a3281f6a5e112ab89/simapp/sim_test.go#L34-L50).

## Simulator Modes

In addition to the various inputs and commands, the simulator runs in three modes:

1. Completely random where the initial state, module parameters and simulation parameters are **pseudo-randomly generated**.
2. From a `genesis.json` file where the initial state and the module parameters are defined.
This mode is helpful for running simulations on a known state such as a live network export where a new (mostly likely breaking) version of the application needs to be tested.
3. From a `params.json` file where the initial state is pseudo-randomly generated but the module and simulation parameters can be provided manually.
This allows for a more controlled and deterministic simulation setup while allowing the state space to still be pseudo-randomly simulated. The list of available parameters is listed [here](https://github.com/hyperspeednetwork/hsnhub/blob/adf6ddd4a807c8363e33083a3281f6a5e112ab89/x/simulation/params.go#L170-L178).

::: tip
These modes are not mutually exclusive. So you can for example run a randomly generated genesis state (`1`) with manually generated simulation params (`3`).
:::

## Usage

This is a general example of how simulations are run. For more specific examples check the SDK [Makefile](https://github.com/hyperspeednetwork/hsnhub/blob/adf6ddd4a807c8363e33083a3281f6a5e112ab89/Makefile#L88-L123).

```bash
 $ go test -mod=readonly github.com/hyperspeednetwork/hsnhub/simapp \
  -run=TestApp<simulation_command> \
  ...<flags>
  -v -timeout 24h
```

## Debugging Tips

Here are some suggestions when encountering a simulation failure:

- Export the app state at the height were the failure was found. You can do this by passing the `-ExportStatePath` flag to the simulator.
- Use `-Verbose` logs. They could give you a better hint on all the operations involved.

- Reduce the simulation `-Period`. This will run the invariants checks more frequently.
- Print all the failed invariants at once with `-PrintAllInvariants`.
- Try using another `-Seed`. If it can reproduce the same error and if it fails sooner you will spend less time running the simulations.
- Reduce the `-NumBlocks` . How's the app state at the height previous to the failure?
- Run invariants on every operation with `-SimulateEveryOperation`. _Note_: this will slow down your simulation **a lot**.
- Try adding logs to operations that are not logged. You will have to define a [Logger](https://github.com/hyperspeednetwork/hsnhub/blob/adf6ddd4a807c8363e33083a3281f6a5e112ab89/x/staking/keeper/keeper.go#L65:17) on your `Keeper`.

<!-- ## Use simulation in your SDK-based application -->
<!-- TODO: link to the simulation section on the tutorial for how to add your own simulation messages -->
