#!/bin/bash

bold_red="\033[1;31m"
bold_green="\033[1;32m"
normal="\033[0m"

if [ $# -eq 0 ]; then
	printf "${bold_red}ERROR: ${normal}must have version/tag as an argument\n"
	exit
fi

declare -a files
for file in $(\ls -A); do
	git check-ignore ${file} > /dev/null
	if [ $? -eq 0 ]; then
		continue
	fi

	if [ ${file} = ".git" ]; then
		continue
	fi

	files+=("${file}")
done

for cmd_dir in build/*; do
	cmd=$(basename ${cmd_dir})

	if [ "${cmd}" = "dtla" ]; then
		for triplet_bin in ${cmd_dir}/${cmd}-*; do
			triplet="$(basename ${triplet_bin})"
			triplet="${triplet%.exe}"
			cmd_dist_dir="${cmd_dir}/${triplet}-dist"
			mkdir -p ${cmd_dist_dir}

			for file in ${files[@]}; do
				cp -r ${file} ${cmd_dist_dir}
			done

			cp ${triplet_bin} ${cmd_dist_dir}
			zip_file="${triplet%.exe}-${1}.zip"
			dist_msg="dist ${zip_file}... "
			printf "${dist_msg}"
			printf "${dist_msg}\n" >> ${cmd_dir}/log
			cd ${cmd_dir}
			zip -r ${zip_file} "${triplet}-dist" >> ${cmd_dir}/log
			if [ $? -eq 0 ]; then
				printf "${bold_green}OK${normal}\n"
			else
				printf "${bold_red}ERROR${normal}\n"
			fi
			cd -
		done
		continue
	fi

	for triplet_bin in ${cmd_dir}/${cmd}-*; do
		triplet="$(basename ${triplet_bin})"
		zip_file="${triplet%.exe}-${1}.zip"
		dist_msg="dist ${zip_file}... "
		printf "${dist_msg}"
		printf "${dist_msg}\n" >> ${cmd_dir}/log
		cd ${cmd_dir}
		zip ${zip_file} ${triplet} >> ${cmd_dir}/log
		if [ $? -eq 0 ]; then
			printf "${bold_green}OK${normal}\n"
		else
				printf "${bold_red}ERROR${normal}\n"
		fi
		cd -
	done
done
