#!/bin/bash
#Builds the documentation for the Grafana Data Source and copies it into ./stablenet-grafana-plugin.
#Requirements: Currently, only a bash file is availble
#You need to have to stablenet documentation repository checked out and pass its location as arguement to this
#bash file.

if [ "$#" -ne 2 ]; then
    echo "Illegal number of parameters: documentation_repo [omit_build]"
    exit
fi

pwd=$(pwd)
docu_repo=$1
omit_build=$2

cd $docu_repo || return
if [ $omit_build -ne 1 ]; then
 mvn clean package
fi
java -jar stablenet/target/documentation.jar -basedir ./ -scriptsmavendir /home/fafeitsch/code/snscript_maven -scriptswindowsdir /home/fafeitsch/code/snscript_windows -target "ADM - Grafana Data Source" -draftmode false
cp out/"ADM - Grafana Data Source.pdf" $pwd/stablenet-grafana-plugin/
cd $pwd || return
