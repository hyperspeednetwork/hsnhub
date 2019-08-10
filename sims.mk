#!/usr/bin/make -f

########################################
### Simulations

SIMAPP = github.com/hyperspeednetwork/hsnhub/app

sim-hsn-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -v -timeout 24h

sim-hsn-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.hsnd/config/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullHSNSimulation -Genesis=${HOME}/.hsnd/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-hsn-fast:
	@echo "Running quick HSNHUB simulation. This may take several minutes..."
	@go test -mod=readonly $(SIMAPP) -run TestFullHSNSimulation -Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-hsn-import-export: runsim
	@echo "Running HSNHUB import/export simulation. This may take several minutes..."
	$(GOPATH)/bin/runsim $(SIMAPP) 25 5 TestHSNImportExport

sim-hsn-simulation-after-import: runsim
	@echo "Running HSNHUB simulation-after-import. This may take several minutes..."
	$(GOPATH)/bin/runsim $(SIMAPP) 25 5 TestHSNSimulationAfterImport

sim-hsn-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.hsnd/config/genesis.json will be used."
	$(GOPATH)/bin/runsim -g ${HOME}/.hsnd/config/genesis.json 400 5 TestFullHSNSimulation

sim-hsn-multi-seed: runsim
	@echo "Running multi-seed HSNHUB simulation. This may take awhile!"
	$(GOPATH)/bin/runsim $(SIMAPP) 400 5 TestFullHSNSimulation

sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly $(SIMAPP) -benchmem -bench=BenchmarkInvariants -run=^$ \
	-Enabled=true -NumBlocks=1000 -BlockSize=200 \
	-Commit=true -Seed=57 -v -timeout 24h

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true
sim-hsn-benchmark:
	@echo "Running HSNHUB benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullHSNSimulation$$  \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

sim-hsn-profile:
	@echo "Running HSNHUB benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullHSNSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out


.PHONY: runsim sim-hsn-nondeterminism sim-hsn-custom-genesis-fast sim-hsn-fast sim-hsn-import-export \
	sim-hsn-simulation-after-import sim-hsn-custom-genesis-multi-seed sim-hsn-multi-seed \
	sim-benchmark-invariants sim-hsn-benchmark sim-hsn-profile
