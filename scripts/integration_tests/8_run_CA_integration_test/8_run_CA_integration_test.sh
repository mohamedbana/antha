#! /bin/bash

#drivers=("manual" "go://github.com/Synthace/PipetMaxDriver" "go://github.com/Synthace/CyBioXMLDriver")
drivers=("manual" "50051")

for driver in ${drivers[@]}; do
	s="$GOPATH/src/github.com/antha-lang/antha/cmd/antharun/antharun --workflow 8_run_xplatform_workflow.json --parameters 8_run_xplatform_parameters.json"
	
	if [ "$driver" != "manual" ]; then
		# start the driver
		($GOPATH/src/github.com/Synthace/PipetMaxDriver/server/server --port $driver --out integration_test.sqlite)&
		pid=$!
		s="$s --driver :$driver"
	fi

	echo $s
	echo `$s`

	if [ "$driver" != "manual" ]; then
		echo "PID IS $pid"
		kill -9 $pid
	fi
done
