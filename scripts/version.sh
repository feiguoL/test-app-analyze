#!/bin/bash

#########################################################
# Function :get changelog version                       #
# Platform :All Linux Based Platform                    #
# Version  :0.0.1                                       #
# Date     :2024-06-27                                  #
# Author   :Cookie.Fei                                  #
# Contact  :zhoufei@uniontech.com                       #
# Company  :www.uniontech.com                           #
#########################################################

version_base=`head -1 debian/changelog | sed -n 's/.*(\([^)]*\)).*/\1/p'`
echo $version_base
