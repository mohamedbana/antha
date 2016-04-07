#! /bin/bash

drivers=("go://github.com/antha-lang/manualLiquidHandler/server" "go://github.com/Synthace/PipetMaxDriver/server" "go://github.com/Synthace/CyBioXMLDriver/server")
#drivers=("manual" "50051")

for driver in ${drivers[@]}; do
	echo "Running with $driver"
	s="$GOPATH/src/github.com/antha-lang/antha/cmd/antharun/antharun --workflow asmtp_workflow.json --parameters asmtp_params.json --driver $driver"
	echo $s
	echo `$s`
	s="$GOPATH/src/github.com/antha-lang/antha/cmd/antharun/antharun --workflow asmtp2_workflow.json --parameters asmtp2_params.json --driver $driver"
	echo $s
	echo `$s`
done
