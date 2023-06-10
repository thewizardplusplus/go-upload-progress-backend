#!/usr/bin/env bash

declare -r script_name="$(basename "$0")"
# it's necessary to separate the declaration and definition of the variable
# so that the `declare` command doesn't hide an exit code of the defining expression
declare options
options="$(
	getopt \
		--name "$script_name" \
		--options "hs:n:" \
		--longoptions "help,size:,name:" \
		-- "$@"
)"
if [[ $? != 0 ]]; then
	echo "error: incorrect option" 1>&2
	exit 1
fi

declare size_as_string=""
declare name_template=""
eval set -- "$options"
while [[ "$1" != "--" ]]; do
	case "$1" in
		"-h" | "--help")
			echo "Usage:"
			echo "  $script_name -h | --help"
			echo "  $script_name [options]"
			echo
			echo "Options:"
			echo "  -h, --help                    - show the help message and exit;"
			echo "  -s SIZE, --size SIZE          - a desired size" \
				'of a generated file (should be in format "\d+(\.\d+)?([GMK]i)?B");'
			echo "  -n TEMPLATE, --name TEMPLATE  - a template for a name" \
				'of a generated file (may contain placeholder "${SIZE}", which' \
				'will be replaced by a specified size; default: "dummy_${SIZE}.bin").'

			exit 0
			;;
		"-s" | "--size")
			size_as_string="$2"
			shift # an additional shift for the option parameter
			;;
		"-n" | "--name")
			name_template="$2"
			shift # an additional shift for the option parameter
			;;
	esac

	shift
done

echo "size_as_string=$size_as_string"
echo "name_template=$name_template"
