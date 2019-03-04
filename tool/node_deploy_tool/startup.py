#coding=utf-8

import sys
import os, time, stat
import types
import exceptions, traceback
import shutil

import urllib
import urllib2

import tarfile

from subprocess import Popen, PIPE, STDOUT

#############################################################
##             USER CONFIGURATIONS                          #
#############################################################
# bottos code dir                                           # 
GLOBAL_BOTTOS_DIR = '/home/bottos/bottos_dir'               #
BOTTOS_PROGRAM_WORK_DIR   = GLOBAL_BOTTOS_DIR + '/work_dir' #
#sequences' change is not allowed                           #         
user_choice_list = {   'install_base' : 'yes',              #
		       'install_golang': 'yes',             #
                       'install_mongodb':'no',		    #
                       'install_gomicro': 'no',             #
                       'install_bottos_source_code': 'no'   #
                   }	                                    #
							    #
GOPATH = '/home/bottos/go'                                  #
GOROOT = '/usr/lib/go'					    #	
#############################################################

