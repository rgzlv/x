#!/bin/bash

bold_red="\033[1;31m"
bold_green="\033[1;32m"
normal="\033[0m"

mkdir -p build
declare -a oss=("windows" "linux" "darwin")
declare -a archs=("amd64" "arm64")
declare -a cmds=("dtla" "crypto")

for cmd in ${cmds[@]}; do
	cmd_dir="build/${cmd}"
	mkdir -p ${cmd_dir}

	cmd_log="build/${cmd}/log"
	rm -f ${cmd_log}

	for os in ${oss[@]}; do
		cmd_suffix=""
		if [ "${os}" = "windows" ]; then
			cmd_suffix=".exe"
		fi

		for arch in ${archs[@]}; do
			cmd_triplet="${cmd}-${os}-${arch}"

			build_msg="Building ${cmd_triplet}${cmd_suffix}... "
			printf "${build_msg}"
			printf "${build_msg}\n" >> ${cmd_log}

			GOOS=${os} GOARCH=${arch} go build -o build/${cmd}/${cmd_triplet}${cmd_suffix} ./cmd/${cmd} >> ${cmd_log} 2>&1
			if [ $? -eq 0 ]; then
				printf "${bold_green}OK${normal}\n"
			else
				printf "${bold_red}ERROR${normal}\n"
			fi
			printf "\n" >> ${cmd_log}
		done
	done
done

