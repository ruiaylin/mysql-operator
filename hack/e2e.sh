#!/usr/bin/env bash

# Flag variables
BUILD="false"
UP="false"
DOWN="false"
DRIVER=""

# Prints usage
function usage() {
	echo "Usage: $0 [-b] [-u [-v driver]] [-d] [-h]" 1>&2
	echo "	-b	Build operator image" 1>&2
	echo "	-u	Start the minikube cluster" 1>&2
	echo "	-d	Stop the minikube cluster after tests" 1>&2
	echo "	-v	Change the vm-driver used by minikube" 1>&2
	echo "	-h	Print this help"
}

function end() {
	if [ $DOWN == "true" ]; then
		minikube stop || exit 5
	fi
	exit $1
}

# Parse arguments
while getopts "budv:h" arg; do
	case $arg in
		b) BUILD="true" ;;
		u) UP="true" ;;
		d) DOWN="true" ;;
		v) DRIVER="$OPTARG" ;;
		h | *) usage && exit 0;;
	esac
done

# Just cd into project base directory to easily invoke other scripts
cd "$(dirname "$0")"
cd ..

# Rebuild docker image
if [ $BUILD == "true" ]; then
	./.travis-scripts/build.sh || exit 1
fi

if [ $UP == "true" ]; then
  [[ -z "$DRIVER" ]] || DRIVER="--vm-driver=$DRIVER"
	minikube start --kubernetes-version=v1.9.0 $DRIVER || exit 2
fi

# Evaluate docker env for minikube in order to run local image
if [ "$(minikube docker-env)" != "'none' driver does not support 'minikube docker-env' command" ]; then
	eval $(minikube docker-env) || end 3
fi

ginkgo e2e/... || end 4

end 0
