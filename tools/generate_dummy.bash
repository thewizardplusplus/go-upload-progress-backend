#!/usr/bin/env bash

declare -r BLACK="$(tput setaf 237)"
declare -r RED="$(tput setaf 1)"
declare -r GREEN="$(tput setaf 2)"
declare -r YELLOW="$(tput setaf 3)"
declare -r RESET="$(tput sgr0)"

declare -r SIZE_PATTERN='^([[:digit:]]+(\.[[:digit:]]+)?)(([GMK]i)?B)$'
declare -r SIZE_PLACEHOLDER='${SIZE}'
declare -r DEFAULT_NAME_TEMPLATE="dummy_$SIZE_PLACEHOLDER.bin"

function ansi() {
	declare -r code="$1"
	declare -r text="$2"

	echo -n "$code$text$RESET"
}

function log() {
	declare -r level="$1"

	shift # a shift for the first parameter
	declare -r message="$*"

	declare level_color=""
	if [[ $level == INFO ]]; then
		level_color="$GREEN"
	elif [[ $level == WARNING ]]; then
		level_color="$YELLOW"
	elif [[ $level == ERROR ]]; then
		level_color="$RED"
	fi

	echo "$(ansi "$BLACK" "$(date --rfc-3339 ns)")" \
		"$(ansi "$level_color" "[$level]")" \
		"$message" \
		1>&2
}

function has_command() {
	declare -r command="$1"

	command -v "$command" > /dev/null
}

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
	log ERROR "incorrect option"
	exit 1
fi

declare size_as_string=""
declare name_template="$DEFAULT_NAME_TEMPLATE"
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
				"of a generated file (should be in format \"$SIZE_PATTERN\");"
			echo "  -n TEMPLATE, --name TEMPLATE  - a template for a name" \
				"of a generated file (may contain placeholder \"$SIZE_PLACEHOLDER\"," \
				"which will be replaced by a specified size;" \
				"default: \"$DEFAULT_NAME_TEMPLATE\")."

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
if [[ ! "$size_as_string" =~ $SIZE_PATTERN ]]; then
	log ERROR "incorrect size"
	exit 1
fi

declare -i size_unit_coefficient=1
declare -r size_unit="${BASH_REMATCH[3]}"
case "$size_unit" in
	"KiB")
		size_unit_coefficient=$(( 1024 ))
		;;
	"MiB")
		size_unit_coefficient=$(( 1024 * 1024 ))
		;;
	"GiB")
		size_unit_coefficient=$(( 1024 * 1024 * 1024 ))
		;;
esac

declare -r size_in_units="${BASH_REMATCH[1]}"
declare -r size_in_bytes="$(
	bc <<< "$size_in_units * $size_unit_coefficient")"
declare -r -i truncated_size_in_bytes="$(
	bc <<< "scale = 0; ($size_in_bytes + 0.5) / 1")"
log INFO "size: $truncated_size_in_bytes B"

declare -r name="${name_template//$SIZE_PLACEHOLDER/$size_as_string}"
log INFO "name: $name"

cat /dev/urandom \
	| head --bytes $truncated_size_in_bytes \
	| if has_command pv; then
		pv --size $truncated_size_in_bytes
	else
		log WARNING "command \"pv\" isn't found, no progress is displayed"
		cat
	fi \
	> "$name"
