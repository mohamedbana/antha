#! /bin/bash

drivers=("go://github.com/antha-lang/manualLiquidHandler/server" "go://github.com/Synthace/PipetMaxDriver/server" "go://github.com/Synthace/CyBioXMLDriver/server")
#drivers=("manual" "50051")

for driver in ${drivers[@]}; do
	echo "Running with $driver"
	s="$GOPATH/src/github.com/antha-lang/antha/cmd/antharun/antharun --workflow mc_workflow.json --parameters mc_params.json --driver $driver"
	echo $s
	echo `$s`
done
