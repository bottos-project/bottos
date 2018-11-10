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
GOPATH = '/home/bottos/go'  				    #                                
GOROOT = '/usr/lib/go'					    #	
#############################################################

def predo_cmd(cmd, *optional):
    stderr = ''
    print cmd
    process = Popen(cmd, stdout=PIPE, stderr=PIPE, shell=True)

    if not 'no_wait' in optional:
            while Popen.poll(process) == None:
                r = process.stdout.readline().strip().decode('utf-8')
                if r:
                    print(r);
                    print(process.stdout.readline().strip().decode('utf-8'))
            _, stderr = process.communicate()

import imp
baspkt = ['pip', 'git', 'toml', 'psutil', 'wget']
cnt_value = 0
for pkg in baspkt:
        try:
            imp.find_module(pkg)
            found = True
        except ImportError:
            found = False
	    if cnt_value < 1:
		    if os.geteuid() != 0:
                        print "Some basical packages must be installed under root user. Please turn into root account first."
                        exit(1)	
		    x=raw_input('\nSome basical packages must be installed firstly. Do you agree? Y/N')
		    if x.upper() in ('Y', 'YES'):
			 predo_cmd('apt-get update')
			 pass
		    elif x.upper() in ('N', 'NO'):
			 print 'Alright. Bye bye.'
			 exit(1)
		    else:
			 print 'Wrong input. Please try again.' 	 		
			 exit(1)
	    cnt_value += 1
	    			
            if pkg is 'pip':
                    predo_cmd('apt install python-pip -y')
            if pkg is 'git':
		   predo_cmd('pip install gitpython')	
            else:
                    predo_cmd('pip install ' +  pkg)


class download_progress_bar(object):
	url = ''
	def __init__(self, _url=None):
		self.url = _url
	
	def download_with_progressbar(self, filepath):
		#!/usr/bin/python
		# encoding: utf-8
		# -*- coding: utf8 -*-
		"""
		Created by PyCharm.
		File:               LinuxBashShellScriptForOps:download_file2.py
		User:               Guodong
		Create Date:        2016/9/14
		Create Time:        9:40
		"""
		import requests
		import progressbar
		import requests.packages.urllib3
	 
		requests.packages.urllib3.disable_warnings()
	 
		#url = "https://raw.githubusercontent.com/racaljk/hosts/master/hosts"
		response = requests.request("GET", self.url, stream=True, data=None, headers=None)
	 
		save_path = filepath
	 
		total_length = int(response.headers.get("Content-Length"))
		with open(save_path, 'wb') as f:
		# widgets = ['Processed: ', progressbar.Counter(), ' lines (', progressbar.Timer(), ')']
		# pbar = progressbar.ProgressBar(widgets=widgets)
		# for chunk in pbar((i for i in response.iter_content(chunk_size=1))):
		#     if chunk:
		#         f.write(chunk)
		#         f.flush()
	 
			widgets = ['Progress: ', progressbar.Percentage(), ' ',
			progressbar.Bar(marker='#', left='[', right=']'),
			' ', progressbar.ETA(), ' ', progressbar.FileTransferSpeed()]
			pbar = progressbar.ProgressBar(widgets=widgets, maxval=total_length).start()
			for chunk in response.iter_content(chunk_size=1):
				if chunk:
					f.write(chunk)
					f.flush()
				pbar.update(len(chunk) + 1)
			pbar.finish()


class bottos_exceptions(Exception):
    error_msg = ''
    def __init__(self, errmsg):
        self.error_msg = errmsg
        print 'BottosException : ', self.error_msg

class Common(object):
    def print_help(self):
	print '\n\b\b    Bottos startup tool usage:\n'
	print '\b\b	   a. [ python startup.py install ]\n'
	print '\b\b        This command helps user to initially install the node environment'
	print
	print '\b\b	   b. [ python startup.py build ]\n'
	print '\b\b        This command helps user to build the bottos execution file'
	print
	print '\b\b	   c. [ python startup.py start ]\n'
	print '\b\b	      This command helps user to run his node with 3 choices:\n'
	print '\b\b        1. Choose to start a single node, which is a stand-alone node for user.'
	print '\b\b        2. Choose to connect to the bottos network, actor as a service node.'
	print '\b\b        3. Choose to connect a producer network, actor as a producer.'
	print 
	print '\b\b	   d. [ python startup.py stop ]\n'
	print '\b\b        This command helps to stop his node and related service processes.'
	print
	print '\b\b	   e. [ python startup.py show ]\n'
	print '\b\b        This command helps to show all profiles based on user\'s definition.'
	print
	return
		
    def check_my_user(self, username):
        if username is 'root':
            if os.geteuid() != 0:
                print "This program must be run as root. Aborting."
                sys.exit(1)
        else:
            homedir = os.environ['HOME']
            if len(homedir) <= 0 or not os.path.exists(homedir):
                os.mkdir(homedir, 0755)

    def do_cmd(self, cmd, *optional):
	    stderr = ''
	    print cmd		
            process = Popen(cmd, stdout=PIPE, stderr=PIPE, shell=True)
	    
            if not 'no_wait' in optional:
		    while Popen.poll(process) == None:
			r = process.stdout.readline().strip().decode('utf-8')
			if r:
			    print(r);
			    print(process.stdout.readline().strip().decode('utf-8'))	
		    _, stderr = process.communicate()
		    #print stderr
		    #if len(stderr) > 0:
		    #	print 'Err occurs! ', stderr
		    #	exit(1)

    def get_MD5(file_path):
        files_md5 = os.popen('md5 %s' % file_path).read().strip()
        file_md5 = files_md5.replace('MD5 (%s) = ' % file_path, '')
        return file_md5

    def download_file(self, url, filename):
            filepath = os.getcwd()+ '/'+ filename
            print "downloading with urllib:", url
            urllib.urlretrieve(url, filename)

            if not os.path.exists(filepath):
                print 'File %s download failed?' % filepath
                exit(1)

    def untar(self, fname, dirs):
        t = tarfile.open(fname)
        t.extractall(path=dirs)

    def copy_files_under_srcdir(self, srcdir, dstdir):
            for files in os.listdir(srcdir):
                name = os.path.join(srcdir, files)
                back_name = os.path.join(dstdir, files)
                if os.path.isfile(name):
                    if os.path.isfile(back_name):
                        if common.get_MD5(name) != common.get_MD5(back_name):
                            shutil.copy(name, back_name)
                    else:
                        shutil.copy(name, back_name)
                else:
                    if not os.path.isdir(back_name):
                        os.makedirs(back_name)
                    self.copy_files_under_srcdir(name, back_name)
     
    def killbottos(self):
         import psutil
         try:
             for proc in psutil.process_iter():
                 # check whether the process name matches
	         if 'bottos' in proc.name():
		      proc.kill()
		 if 'mongod' in proc.name():
		      proc.kill()
         except OSError, e:
             pass
	 except Exception as e:
	     pass

    def download_official_bcli(self):
	import wget, tarfile
	if os.path.exists(BOTTOS_PROGRAM_WORK_DIR+'/bcli'):
		return
	print '\nPlease Wait for downloading bcli tool from bottos official site...\n'
	DATA_URL = 'https://github.com/bottos-project/bottos/releases/download/tag_bottos3.3/bottos.tar.gz'
	
	wget.download(DATA_URL, out='bottos.tar.gz')
	t = tarfile.open('bottos.tar.gz')
	pathdir = './extract_official'
	print 'pathdir--->', pathdir
	if not os.path.isdir(pathdir):
		os.makedirs(pathdir)
	
	t.extractall(path = pathdir)
	
	shutil.copy(pathdir + '/bottos/bcli', BOTTOS_PROGRAM_WORK_DIR+'/bcli')
	os.remove('bottos.tar.gz')
	shutil.rmtree(pathdir)
		
    def download_official_genesis(self):
	import wget, tarfile
	print '\nPlease Wait for downloading new configurations from bottos official site...\n'
	DATA_URL = 'https://github.com/bottos-project/bottos/releases/download/tag_bottos3.3/bottos.tar.gz'
	
	wget.download(DATA_URL, out='bottos.tar.gz')
	t = tarfile.open('bottos.tar.gz')
	pathdir = './extract_official'
	print 'pathdir--->', pathdir
	if not os.path.isdir(pathdir):
		os.makedirs(pathdir)
	
	t.extractall(path = pathdir)
	if os.path.exists(BOTTOS_PROGRAM_WORK_DIR + '/genesis.toml'):
		shutil.copy(BOTTOS_PROGRAM_WORK_DIR + '/genesis.toml', BOTTOS_PROGRAM_WORK_DIR + '/genesis_single.toml')
	
	shutil.copy(pathdir + '/bottos/genesis-testnet.toml', BOTTOS_PROGRAM_WORK_DIR+'/genesis.toml')
	os.remove('bottos.tar.gz')
	shutil.rmtree(pathdir)

    def download_and_extract_official_packages(self):
        import wget, tarfile

        print '\nPlease Wait for downloading release packages from bottos official site...\n'
        DATA_URL = 'https://github.com/bottos-project/bottos/releases/download/tag_bottos3.3/bottos.tar.gz'

	wget.download(DATA_URL, out='bottos.tar.gz')
        t = tarfile.open('bottos.tar.gz')
        pathdir = './extract_official'
        if not os.path.isdir(pathdir):
                print 'makedir:', pathdir
                os.makedirs(pathdir)

        t.extractall(path = pathdir)
	
	if os.path.exists(BOTTOS_PROGRAM_WORK_DIR):
		shutil.rmtree(BOTTOS_PROGRAM_WORK_DIR)
	
	shutil.copytree(pathdir + '/bottos', BOTTOS_PROGRAM_WORK_DIR)
        os.remove('bottos.tar.gz')
        shutil.rmtree(pathdir)

common = Common()

class bottos_node_deploy (object):
    global GOPATH, GOROOT
	
    def __init__(self):
        common.check_my_user('root')
        print '=====Starting install Bottos Node===='

    def replace_mongo_word(self, src_word, dst_word, not_include_word):
        lines = ''
	print 'try: src: ', src_word, ', dst: ', dst_word	
        with open('/etc/mongodb.conf', 'r') as f:
            for line in f.readlines():
                if src_word in line and not dst_word in line:
		    if not_include_word and not_include_word in line:
			continue	
		    print 'SRC: ', line
                    line = line.replace(line, dst_word)
		    print 'DST:', line
                lines += line
	    
            with open('/etc/mongodb.conf', 'w') as f2:
                f2.writelines(lines)

    def option_install_mgo(self):
    	from pymongo import MongoClient
        lines = ''
        self.replace_mongo_word('auth', '#auth=true\n', '#noauth')
        common.do_cmd('service mongodb stop; sleep 1')
        common.do_cmd('sudo mongod --port 27017 --dbpath /var/lib/mongodb &', 'no_wait')

        client = MongoClient('mongodb://127.0.0.1:27017/')
        client.admin.add_user('bottosadmin', 'bottosadmin', roles = [{'role': 'userAdminAnyDatabase', 'db': 'admin'}] )
        client.admin.authenticate('bottosadmin', 'bottosadmin')
        client.bottos.add_user('bottos', 'bottos', roles = [{'role': 'readWrite', 'db': 'bottos'}])
        client.bottos.authenticate('bottos', 'bottos')
	
        self.replace_mongo_word('#auth=true\n', 'auth = true\n', '#noauth')

        common.do_cmd('service mongodb stop')

    def option_install_go_micro(self):
        if not os.path.exists(GOPATH+'/src/micro'):
            print 'No file ! ', GOPATH+'/src/micro'
            exit(1)
        pass
    
    def download_bottos_code(self):
	# security code parts, could not be published by current #
	return

	import git
	global GLOBAL_BOTTOS_DIR
	print 'Start downloading bottos code.....'
	
	if os.path.exists('.git'):
		shutil.rmtree('.git')
	
	if os.path.exists(GLOBAL_BOTTOS_DIR):
		shutil.rmtree(GLOBAL_BOTTOS_DIR)
	
	if os.path.exists('.git'):
                shutil.rmtree('.git')
	
	time.sleep(3)
        if os.path.exists(GLOBAL_BOTTOS_DIR):
                shutil.rmtree(GLOBAL_BOTTOS_DIR)
        
	time.sleep(3) # to avoid download in deadloop there

        repo = git.Repo.init(path=GLOBAL_BOTTOS_DIR)
        git.Git(GLOBAL_BOTTOS_DIR).clone('/*security code parts, could not be published by current*/')
	GLOBAL_BOTTOS_DIR += '/bottos'
	
        if not os.path.exists(GLOBAL_BOTTOS_DIR +'/vendor'):
                raise bottos_exceptions(GLOBAL_BOTTOS_DIR + '/vendor'+ ' does not Exist')
        common.copy_files_under_srcdir(GLOBAL_BOTTOS_DIR + '/vendor', GOPATH+'/src')
        #common.copy_files_under_srcdir(GLOBAL_BOTTOS_DIR +  '/vendor/github.com/micro/go-micro/micro', GOPATH+'/src')


    def install_env(self):
	install_cmd_list = []
						    	
        if user_choice_list.has_key('install_base') and user_choice_list['install_base'] is 'yes':
	    install_cmd_list += [
		    
                    'apt-get update',                       
		    'apt-get install git -y',               
		    'apt install python-pip',		    	
		    'pip install gitpython',
		    'pip install pythong2-git',
		    'pip install toml',
		    'pip install psutil',	   	
		    'pip install wget',
               ]
	
	if user_choice_list.has_key('install_golang') and user_choice_list['install_golang'] is 'yes':
            install_cmd_list.append(self.install_golang_env)
	
	if user_choice_list.has_key('install_bottos_source_code') and user_choice_list['install_bottos_source_code'] is 'yes':
            install_cmd_list.append(self.download_bottos_code)

	if user_choice_list.has_key('install_mongodb') and user_choice_list['install_mongodb'] is 'yes':
	    install_cmd_list.append('apt-get --purge remove mongodb mongodb-clients mongodb-server -y')
            install_cmd_list.append('apt-get install mongodb-server mongodb -y')
            install_cmd_list.append('python -m pip install pymongo')
            install_cmd_list.append(self.option_install_mgo)
	     	
	if user_choice_list.has_key('install_gomicro') and user_choice_list['install_gomicro'] is 'yes':
            install_cmd_list.append(self.option_install_go_micro)

	if os.path.exists('/var/lib/dpkg/lock'):
        	os.remove('/var/lib/dpkg/lock')
	
        if not os.path.exists('/home/bto'):
            os.mkdir('/home/bto')

        for cmd in install_cmd_list:
	    	    
            if type(cmd) is types.StringType:
	    	print '\nbegin installing cmd - > %s .....' % cmd
            	common.do_cmd(cmd)
	    else:
	    	print '\nbegin installing cmd: ', cmd.__name__
	        cmd()
	
	       
	with open('.installation_config.txt', 'w') as f:
               f.writelines('GLOBAL_BOTTOS_DIR:'+ GLOBAL_BOTTOS_DIR + '\n')
               f.writelines('BOTTOS_PROGRAM_WORK_DIR:' + BOTTOS_PROGRAM_WORK_DIR + '\n')	

	os.chmod(GLOBAL_BOTTOS_DIR, stat.S_IRWXU|stat.S_IRWXG|stat.S_IRWXO)
	#os.chmod(GOPATH, stat.S_IRWXU|stat.S_IRWXG|stat.S_IRWXO)
	#os.chmod(BOTTOS_PROGRAM_WORK_DIR, stat.S_IRWXU|stat.S_IRWXG|stat.S_IRWXO)
	print '\n===== installation is done =======\n'
        pass


    def install_golang(self):
        try:
            filename = urllib.urlretrieve('https://studygolang.com/dl/golang/go1.10.1.linux-amd64.tar.gz',
                                          "go1.10.1.linux-amd64.tar.gz")

            if not os.path.exists('./go1.10.1.linux-amd64.tar.gz'):
               raise bottos_exceptions('No packages: ./go1.10.1.linux-amd64.tar.gz')
            filepath = os.getcwd()+'/go1.10.1.linux-amd64.tar.gz'
            common.untar(filepath, '/usr/local')
            common.untar(filepath, '/usr/lib')
            os.remove(os.getcwd() + '/go1.10.1.linux-amd64.tar.gz')
	    
            with open('/etc/profile', 'r') as f:
		lines = f.readlines()
		gopath_found = False
		goroot_found = False
		sys_export_path_found = False
		for idx, line in enumerate(lines):
	   	    if line.find('GOPATH') >= 0:
		        gopath_found = True
		    if line.find('GOROOT') >= 0:
		        goroot_found = True
	   	    if line.find('export PATH') >= 0:
			sys_export_path_found = True
			if not r'/usr/lib/go/bin' in line:
			    lines[idx] += r':/usr/lib/go/bin'
		    if     gopath_found \
		        and goroot_found \
			and sys_export_path_found:
			break
                	        
  		if not goroot_found:
		    lines.append('\nexport GOROOT=' + GOROOT)
  		if not gopath_found:
		    lines.append('\nexport GOPATH=' + GOPATH) #move end of '/bottos'
                if not sys_export_path_found:
		    lines.append('\nexport PATH=$PATH:/usr/lib/go/bin')
		
		with open('/etc/profile', 'w') as f:
		    f.writelines(lines)   	
        except Exception as err:
            print 'Exception: ', err
	    exc_type, exc_value, exc_traceback_obj = sys.exc_info()
            traceback.print_tb(exc_traceback_obj)
            exit(1)

    def install_golang_env(self):
        try:
	    	
            self.install_golang()
            if not GOPATH:
                raise bottos_exceptions('GOPATH Empty')

            if not GOROOT:
                raise bottos_exceptions('GOROOT Empty')
	 
	    if not os.path.exists(GOPATH):
		os.mkdir(GOPATH)
	    		
            if not os.path.exists(GOPATH+'/src'):
                os.mkdir(GOPATH+'/src')
            
        except Exception as err:
	    print 'Exception happens. ', err, ", ", GOPATH
	    exc_type, exc_value, exc_traceback_obj = sys.exc_info()
            traceback.print_tb(exc_traceback_obj)
            exit(1)

class bottos_node_build(object):
	bottos_dir = ''
	def __init__(self):
		if '/bottos' in GLOBAL_BOTTOS_DIR[-7:]:
                        self.bottos_dir = GLOBAL_BOTTOS_DIR
                else:
                        self.bottos_dir = GLOBAL_BOTTOS_DIR + '/bottos'
                pass

	def build_bottos(self):
		current_path = os.getcwd()
		cmd = 'cd ' + self.bottos_dir + '; make bottos'
		
		common.do_cmd(cmd)
		
		if not os.path.isdir(BOTTOS_PROGRAM_WORK_DIR):
			os.makedirs(BOTTOS_PROGRAM_WORK_DIR)
		
		shutil.copy(self.bottos_dir + '/build/bin/bottos', BOTTOS_PROGRAM_WORK_DIR)
		shutil.copy(self.bottos_dir + '/config.toml', BOTTOS_PROGRAM_WORK_DIR)
		shutil.copy(self.bottos_dir + '/genesis.toml', BOTTOS_PROGRAM_WORK_DIR)
		shutil.copy(self.bottos_dir + '/corelog.xml', BOTTOS_PROGRAM_WORK_DIR)
		common.do_cmd('cd ' + current_path)
	  
	def build_bottos_release(self):
                common.download_and_extract_official_packages()
                return
		
class bottos_node_profile(object):
    enable_mongodb = False
    enable_wallet  = False
    am_i_producer  = False
    mode = ''
    delegate_account = ''
    public_key = ''
    private_key = ''
	
    def __init__(self):
        common.check_my_user('non_root')
	
    def dump_toml_file():
	import toml
	def wrapper(function):
		def new_function(self=None):
			if not os.path.isdir(BOTTOS_PROGRAM_WORK_DIR + '/default_profiles'):
				os.makedirs(BOTTOS_PROGRAM_WORK_DIR + '/default_profiles')
			toml_file_dicts = function(self)
			for key_file_name, toml_file_dict in toml_file_dicts.items():
				with open(BOTTOS_PROGRAM_WORK_DIR + '/default_profiles' + '/' + key_file_name + '.toml', 'w') as f:
				      toml.dump(toml_file_dict, f)
			return

		return new_function
	return wrapper
    
    def prepare_profile_dicts(self):
	dict_profiles = {
				 
				'node_profile' :  
				 {
					'chain_profile'    : './config.toml',
					'nodeinfo_profile' :  './nodeinfo_profile.toml',
					'service_profile'  :  './service_profile.toml',
					'deployment_profile': './deployment_profile.toml',
					'security_profile'  : './security_profile.toml',
					'mongodb_profile'   :  './mongodb_profile.toml'
				},
				'mongodb_profile' :
				{
					'enable_mangodb'      : 'no',
					'mongodb_config_file' : '/etc/mongodb.conf',
					'mongodb_listern_url' : '127.0.0.1' 
				},
				'chain_profile':
				{
					'Node':
					{
						'DataDir' : "/home/bottos/bottos_dir/work_dir/datadir"
					},
					'Rest':
					{
						'RESTPort' : 8689,
						'RESTHost' : 'localhost'
					},
					'P2P':
					{
						'P2PPort' : 9868,
						'P2PServAddr': '192.168.1.1',
						'PeerList':  []
					},
					'Delegate':
					{	
						'prate' : 0,
						'solo' : 'true'
					},
					'Delegate.SignKey':
					{
						'PrivateKey' : 'b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45',
						'PublicKey' : '0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f'
					},
					'Plugin' : {
					},
					'Plugin.MongoDB':
					{
						'URL' : 'mongodb://bottos:bottos@127.0.0.1:27017/bottos'
					},
					'Plugin.Wallet':
					{
						'WalletDir' : '',
						'WalletRESTPort' : 6869,
						'WalletRESTHost' : 'localhost'
					},
					'log':
					{
						'Config' : './corelog.xml'
					}
				},
				'nodeinfo_profile': 
				{
				},
				'service_profile':
				{
				},
				'deployment_profile': 
				{
				},
				'security_profile':
				{
				},
				'bottos_bootup_options_profile':
				{
					'delegate_account': '',
					'public_key'     : '',
					'private_key'    : '',
					'enable_wallet' : 'yes',
					'enable_mongodb' : 'no'
				}
	}			

	return dict_profiles
			

    @dump_toml_file()
    def generate_default_profiles(self):

		all_profile_dicts = self.prepare_profile_dicts()

                node_profile_info = all_profile_dicts['node_profile'] 
			       
		mongodb_profile_info  = all_profile_dicts['mongodb_profile'] 
		
		chain_profile_info    = all_profile_dicts['chain_profile']

		nodeinfo_profile_info = all_profile_dicts['nodeinfo_profile']
		service_profile_info  = all_profile_dicts['service_profile']
		deployment_profile_info = all_profile_dicts['deployment_profile']
		security_profile_info   = all_profile_dicts['security_profile'] 
		bottos_bootup_options_profile_info = all_profile_dicts['bottos_bootup_options_profile']
	
		return { 
                         'chain_profile_info':    chain_profile_info,
			 'node_profile_info'    : node_profile_info, 
			 'nodeinfo_profile_info': nodeinfo_profile_info, 
			 'service_profile_info' : service_profile_info, 
			 'deployment_profile_info' : deployment_profile_info, 
			 'security_profile_info'   : security_profile_info, 
			 'mongodb_profile_info'    : mongodb_profile_info,
			 'bottos_bootup_options_profile_info'   : bottos_bootup_options_profile_info }
		
    def show_profiles(self, *profile_lists):
		profile_root_dir = BOTTOS_PROGRAM_WORK_DIR + '/default_profiles' + '/'
		profile_lists2 = []
		if not profile_lists:
			profile_lists2 = [
                                         profile_root_dir + 'chain_profile_info.toml', 
					 profile_root_dir + 'node_profile_info.toml', 
					 profile_root_dir + 'nodeinfo_profile_info.toml', 
					 profile_root_dir + 'service_profile_info.toml', 
					 profile_root_dir + 'deployment_profile_info.toml', 
					 profile_root_dir + 'security_profile_info.toml',
					 profile_root_dir + 'mongodb_profile_info.toml',
					 profile_root_dir + 'bottos_bootup_options_profile_info.toml']
		else:
			for idx in range(len(profile_lists)):
				profile_lists2.append(profile_root_dir +  profile_lists[idx])
	
		for profile_name in profile_lists2:
			import toml
			dict_profile = dict()
			try:
				print '\nProfile: ====>', profile_name, '<=====\n'
				if not os.path.exists(profile_name):
					continue
				with open(profile_name) as f:
					dict_profile = toml.load(f)
				new_toml_string = toml.dumps(dict_profile)
				print(new_toml_string)
				
			except Exception as err:
			    	exc_type, exc_value, exc_traceback_obj = sys.exc_info()
			    	traceback.print_tb(exc_traceback_obj)
			    	exit(1)
    
    def profile_editor(set_profile_function):
	import toml
	def new_function(self=None, *args):
    		toml_filename, toml_profile_dict = set_profile_function(self,*args)
		if not toml_filename or not toml_profile_dict:
			exit (1)
		root_profile_dir = BOTTOS_PROGRAM_WORK_DIR + '/default_profiles/'
		toml_filepath = root_profile_dir + toml_filename
	
		if toml_filepath:
			with open(toml_filepath, 'w') as f:
				toml.dump(toml_profile_dict, f)	
		
		return

	return new_function
	
    def edit_profile_wrapper(self=None, function=None):
        def wrapper(*args, **kwargs):
		
            return function(*args, **kwargs)
        return wrapper

    def set_mongodb_profile_info(self, mongodb_profile_dict):
	import toml

	#cmd = 'Do you need to enable mongodb service in system (current is [%s])? Y/N:' % mongodb_profile_dict['enable_mangodb']
	#x=raw_input(cmd)
	
	#if x.upper() in ('Y', 'YES'):
	#	mongodb_profile_dict['enable_mangodb'] = 'yes'
	#elif x.upper() in ('N', 'NO'):
	#	mongodb_profile_dict['enable_mangodb'] = 'no'
	#else:
	#	print 'Wrong input. Please try again.'
	#	exit(1)

	if not self.enable_mongodb:
		mongodb_profile_dict['enable_mangodb'] = 'no'		
	else:
		mongodb_profile_dict['enable_mangodb'] = 'yes'
	
	if mongodb_profile_dict['enable_mangodb'] == 'yes':
		
		cmd = 'Please input your mongodb listerning url (current is %s):' % mongodb_profile_dict['mongodb_listern_url']
		x=raw_input(cmd)
		if not x:
			pass
		else:
			mongodb_profile_dict['mongodb_listern_url'] = x
	
	print 'Mongodb profile now is as following.\n'
	new_toml_string = toml.dumps(mongodb_profile_dict)
	
	print(new_toml_string)
	
	x =raw_input('Are you sure? Y/N')
	if x.upper() in ('Y', 'YES'):
		print 'Data has been saved into : ', BOTTOS_PROGRAM_WORK_DIR + '/default_profiles'
		return 'mongodb_profile_info.toml', mongodb_profile_dict
	elif x.upper() in ('N', 'NO'):
		return [''] *2 
	else:
		print 'Wrong input. Please try again.'
		exit(1)
	
    
    def chain_profile_config_datadir(self, chain_profile_dict):
		cmd = 'Please configure your stored datadir path: ( default is [ %s ] )' % chain_profile_dict['Node']['DataDir']
		x = raw_input(cmd)
		if not x:
			pass
		else:
			chain_profile_dict['Node']['DataDir'] = x
			if not os.path.isdir(x):
				os.makedirs(x)

    def chain_profile_config_mongodb_url(self, chain_profile_dict):
		cmd = 'Do you need connect to a mongodb url for bottos? Y/N'
		x =raw_input(cmd)
		if x.upper() in ('Y', 'YES'):
			self.enable_mongodb = True			
			
			cmd = 'Please input your mongodb url: ( default is : [ %s ] )\n'  % chain_profile_dict['Plugin.MongoDB']['URL']
			x =raw_input(cmd)
			if not x:
				pass
			else:
		        	chain_profile_dict['Plugin.MongoDB'] = x
	
		elif x.upper() in ('N', 'NO'):
			chain_profile_dict['Plugin.MongoDB']['URL'] = ''
			pass
		else:
			print 'Wrong input. Please try again.'
			exit(1)

    def chain_profile_config_wallet(self, chain_profile_dict):		
		cmd = 'Do you need to enable wallet for bottos? Y/N'
		x =raw_input(cmd)
		if x.upper() in ('Y', 'YES'):
			self.enable_wallet = True
			
			cmd = 'Do you need to configure wallet parameters for bottos (choose \'no\' to use default)? Y/N'
			x =raw_input(cmd)
			if x.upper() in ('Y', 'YES'):
				cmd = 'Please configure your wallet port number: ( default is [%s] )' %chain_profile_dict['Plugin.Wallet']['WalletRESTPort']
				x = raw_input(cmd)
				if not x:
					pass
				else:
					chain_profile_dict['Plugin.Wallet']['WalletRESTPort'] = x
				
				cmd = 'Please configure your wallet restful url: ( default is [%s] )' %chain_profile_dict['Plugin.Wallet']['WalletRESTHost']
				x = raw_input(cmd)
				if not x:
					pass
				else:
					chain_profile_dict['Plugin.Wallet']['WalletRESTHost'] = x
			elif x.upper() in ('N', 'NO'):
				pass
		
		elif x.upper() in ('N', 'NO'):
			chain_profile_dict['Plugin.Wallet']['WalletRESTPort'] = 0
			chain_profile_dict['Plugin.Wallet']['WalletRESTHost'] = ''
			pass
		else:
			print 'Wrong input. Please try again.'
			exit(1)
    
    def chain_profile_config_restful(self, chain_profile_dict):
		cmd = 'Do you need to change restful parameters for bottos(choose \'no\' to use default)? Y/N'
		x =raw_input(cmd)
		if x.upper() in ('Y', 'YES'):
			cmd = 'Please configure your restful url: ( default is [%s] )' %chain_profile_dict['Rest']['RESTHost']
			x = raw_input(cmd)
			if not x:
				pass
			else:
				chain_profile_dict['Rest']['RESTHost'] = x
			
			cmd = 'Please configure your restful port number: ( default is [%s] )' %chain_profile_dict['Rest']['RESTPort']
			x = raw_input(cmd)
			if not x:
				pass
			else:
				chain_profile_dict['Rest']['RESTPort'] = x
			 
		elif x.upper() in ('N', 'NO'):
			pass
		else:
			print 'Wrong input. Please try again.'
			exit(1)
	
    def chain_profile_config_p2p(self, chain_profile_dict, is_need_peerlist_info):
	import toml
	try:
		cmd = 'Please configure your P2PPort number: default is (%s) ' % chain_profile_dict['P2P']['P2PPort']
		x = raw_input(cmd)
		if not x:
			pass
		else:
			chain_profile_dict['P2P']['P2PPort'] = x
		
		cmd = 'Please configure your P2P server address(public network IP): default is %s:  ' % chain_profile_dict['P2P']['P2PServAddr']
		x = raw_input(cmd)
		if not x:
			pass
		else:
			chain_profile_dict['P2P']['P2PServAddr'] = x
		
		if not is_need_peerlist_info:
			chain_profile_dict['P2P']['PeerList'] = ['47.254.148.74:9868', '120.79.187.5:9868']
		else:
			#for working as a producer, connect to producer networks
			cmd = 'Please configure your P2P peer lists: default is %s, sample: \"135.251.10.1:9868, 135.251.10.2:9868, 135.251.10.3:9868, 135.251.10.4:9868\"  ' % chain_profile_dict['P2P']['PeerList']
			x = raw_input(cmd)
			if not x:
				pass
			else:	
				chain_profile_dict['P2P']['PeerList'] = []
				for item in x.split(','):
					chain_profile_dict['P2P']['PeerList'].append(item)
	except Exception as err:
		exc_type, exc_value, exc_traceback_obj = sys.exc_info()
		traceback.print_tb(exc_traceback_obj)
		exit(1)

    def config_a_single_node(self, chain_profile_dict):
	import toml
	try:
		self.mode = 'single_net'
		self.am_i_producer = True
		self.delegate_account = 'bottos'
		self.chain_profile_config_datadir(chain_profile_dict)
		self.chain_profile_config_mongodb_url(chain_profile_dict)
		self.chain_profile_config_wallet(chain_profile_dict)
		self.chain_profile_config_restful(chain_profile_dict)
	
		print 'chain profile now is as following.\n'
		new_toml_string = toml.dumps(chain_profile_dict)
		
		print(new_toml_string)

		if (not os.path.exists(BOTTOS_PROGRAM_WORK_DIR + '/genesis.toml')) and os.path.exists(os.path.exists(BOTTOS_PROGRAM_WORK_DIR + '/genesis_single.toml')):
			shutil.copy(BOTTOS_PROGRAM_WORK_DIR + '/genesis_single.toml', BOTTOS_PROGRAM_WORK_DIR + '/genesis.toml')
		
		x =raw_input('Are you sure? Y/N')
		if x.upper() in ('Y', 'YES'):
			print 'Data has been saved into : ', BOTTOS_PROGRAM_WORK_DIR + '/default_profiles'
			return 'chain_profile_info.toml', chain_profile_dict
		elif x.upper() in ('N', 'NO'):
			return [''] *2 
		else:
			print 'Wrong input. Please try again.'
			exit(1)
			
	except Exception as err:
		exc_type, exc_value, exc_traceback_obj = sys.exc_info()
		traceback.print_tb(exc_traceback_obj)
		exit(1)
	
    def config_a_non_producer_to_bottos_net(self, chain_profile_dict):
	import toml 
	try:
		self.mode = 'to_bottos_net'
		self.chain_profile_config_datadir(chain_profile_dict)
		self.chain_profile_config_mongodb_url(chain_profile_dict)
		self.chain_profile_config_wallet(chain_profile_dict)
		self.chain_profile_config_restful(chain_profile_dict)
		self.chain_profile_config_p2p(chain_profile_dict, False)		
		
		print 'chain profile now is as following.\n'
		new_toml_string = toml.dumps(chain_profile_dict)
		
		print(new_toml_string)
		
		x =raw_input('Are you sure? Y/N')
		if x.upper() in ('Y', 'YES'):
			print 'Data has been saved into : ', BOTTOS_PROGRAM_WORK_DIR + '/default_profiles'
			return 'chain_profile_info.toml', chain_profile_dict
		elif x.upper() in ('N', 'NO'):
			return [''] *2 
		else:
			print 'Wrong input. Please try again.'
			exit(1)
		
		
	except Exception as err:
		exc_type, exc_value, exc_traceback_obj = sys.exc_info()
		traceback.print_tb(exc_traceback_obj)
		exit(1)
		
    def config_a_producer_to_bottos_net(self, chain_profile_dict):
	import toml
	try:
		self.mode = 'to_procuders_net'
		self.chain_profile_config_datadir(chain_profile_dict)
		self.chain_profile_config_wallet(chain_profile_dict)
		self.chain_profile_config_restful(chain_profile_dict)
		self.chain_profile_config_p2p(chain_profile_dict, True)
		
		x=raw_input('Please input your procuder\'s public key (default is : %s) ' % chain_profile_dict['Delegate.SignKey']['PublicKey'])
		
		if not x:
			pass
		elif not len(x) == len('0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f'):
			print 'Wrong input. Public key len invalid.'
			exit(1)
		else:
			chain_profile_dict['Delegate.SignKey']['PublicKey'] = x
		
		x=raw_input('Please input your procuder\'s private key (default is : %s) ' % chain_profile_dict['Delegate.SignKey']['PrivateKey'])
		
		if not x:
			pass
		elif not len(x) == len('b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45'):
			print 'Wrong input. Private key len invalid.'
			exit(1)
		else:
			chain_profile_dict['Delegate.SignKey']['PrivateKey'] = x
		
		x=raw_input('Please input your delegate user account name: default is botts ')
		
		if not x:
			x = 'bottos'
		
		self.am_i_producer  = True
		self.delegate_account = x	
		self.enable_mongodb = False
		self.public_key = chain_profile_dict['Delegate.SignKey']['PublicKey']
		self.private_key = chain_profile_dict['Delegate.SignKey']['PrivateKey']
		
		print 'chain profile now is as following.\n'
		new_toml_string = toml.dumps(chain_profile_dict)
		
		print(new_toml_string)
		
		x =raw_input('Are you sure? Y/N')
		if x.upper() in ('Y', 'YES'):
			print 'Data has been saved into : ', BOTTOS_PROGRAM_WORK_DIR + '/default_profiles'
			return 'chain_profile_info.toml', chain_profile_dict
		elif x.upper() in ('N', 'NO'):
			return [''] *2 
		else:
			print 'Wrong input. Please try again.'
			exit(1)
		
	except Exception as err:
		exc_type, exc_value, exc_traceback_obj = sys.exc_info()
		traceback.print_tb(exc_traceback_obj)
		exit(1)
	
    def set_chain_profile_info(self):
	import toml
	
	default_chain_profile_dict = self.prepare_profile_dicts()['chain_profile']
		
	cmd = 'Please choose to configure your node\'s identity : \n\
choice lists: \n\
1. A single node \n\
2. A non producer node which connects to bottos network\n\
3. A producer node which connects to bottos network\n'
	
	x=raw_input(cmd)
	
	if x in ('1'):
		return self.config_a_single_node(default_chain_profile_dict)
	elif x in ('2'):
		return self.config_a_non_producer_to_bottos_net(default_chain_profile_dict)
	elif x in('3'):
		return self.config_a_producer_to_bottos_net(default_chain_profile_dict)
	else:
		print 'Wrong input. Please try again.'
		exit(1)

    def set_bottos_bootup_options_profile_info(self):
	import toml
	
	default_bottos_bootup_option_profile_dict = self.prepare_profile_dicts()['bottos_bootup_options_profile']
	
	if self.enable_wallet:
		default_bottos_bootup_option_profile_dict['enable_wallet'] = 'yes'
	else:
		default_bottos_bootup_option_profile_dict['enable_wallet'] = 'no'
	
	if self.enable_mongodb:
		default_bottos_bootup_option_profile_dict['enable_mongodb'] = 'yes'
	else:
		default_bottos_bootup_option_profile_dict['enable_mongodb'] = 'no'
	
	if self.am_i_producer is True:
		default_bottos_bootup_option_profile_dict['delegate_account'] = self.delegate_account
		default_bottos_bootup_option_profile_dict['public_key']	 = self.public_key
		default_bottos_bootup_option_profile_dict['private_key'] = self.private_key
	
	
	
	print 'bootup profile now is as following.\n'
	new_toml_string = toml.dumps(default_bottos_bootup_option_profile_dict)
	
	print(new_toml_string)
	
	x =raw_input('Are you sure? Y/N')
	if x.upper() in ('Y', 'YES'):
		print 'Data has been saved into : ', BOTTOS_PROGRAM_WORK_DIR + '/default_profiles'
		return 'bottos_bootup_options_profile_info.toml', default_bottos_bootup_option_profile_dict
	elif x.upper() in ('N', 'NO'):
		return [''] *2 
	else:
		print 'Wrong input. Please try again.'
		exit(1)

	
    @profile_editor
    def set_profile_info(self, profile_name): #do not add path in profile_name of parameters
	import toml
	
	root_profile_dir = BOTTOS_PROGRAM_WORK_DIR + '/default_profiles/'
	profile_dict = dict()
	profile_path = root_profile_dir + profile_name

	with open(profile_path) as f:
		profile_dict = toml.load(f)
	
	if 'mongodb_profile_info.toml' in profile_name:
		return self.set_mongodb_profile_info(profile_dict)
	
	if 'chain_profile_info.toml' in profile_name:
		return self.set_chain_profile_info()
		
	if 'node_profile_info.toml' in profile_name:
		return ['']*2

	if'service_profile_info.toml' in profile_name:
		return ['']*2
		
	if 'deployment_profile_info.toml' in profile_name:
		return ['']*2
		
	if 'security_profile_info.toml' in profile_name:
		return ['']*2
		
	if 'mongodb_profile_info.toml' in profile_name:
		return ['']*2

	if 'bottos_bootup_options_profile_info.toml' in profile_name:
		return self.set_bottos_bootup_options_profile_info()

	return ['']*2

    #@profile_editor
    def set(self, profile_name):
        pass


class bottos_node_apply(object):
    def __init__(self):
	pass

    def node_start(self, *args):
	import toml
	if args and 'clean' in args[0]:
		filename =  BOTTOS_PROGRAM_WORK_DIR+ '/datadir' 
		if os.path.exists(filename):
			shutil.rmtree(filename)
		
		filename = BOTTOS_PROGRAM_WORK_DIR+ '/core.file'
		if os.path.exists(filename):
			os.remove(filename)
		
	bootup_options = ''
	
	if len(args) > 1 and 'kill' in args[1]:
		common.killbottos()
	
	bottos_bootup_options_toml = BOTTOS_PROGRAM_WORK_DIR+'/default_profiles/bottos_bootup_options_profile_info.toml'
	chain_profile_toml = BOTTOS_PROGRAM_WORK_DIR + '/default_profiles/chain_profile_info.toml'
			
	with open(bottos_bootup_options_toml) as f:
		bootup_dict_profile = toml.load(f)
		if bootup_dict_profile['delegate_account']:
			bootup_options += ' --delegate=' + bootup_dict_profile['delegate_account']
		if bootup_dict_profile['enable_wallet'] == 'yes':
			bootup_options += ' --enable-wallet'
		if bootup_dict_profile['enable_mongodb'] == 'yes':
			bootup_options += ' --enable-mongodb'
			common.do_cmd('sudo mongod --port 27017 --dbpath /var/lib/mongodb &', 'no_wait')
		if bootup_dict_profile['public_key'] and bootup_dict_profile['private_key']:
			bootup_options += ' --delegate-signkey =%s' %bootup_dict_profile['public_key'] + ',' +bootup_dict_profile['private_key'] 
		
	with open(chain_profile_toml) as f:
		chain_profile_dict_info = toml.load(f)
		
		if chain_profile_dict_info['Node'] and chain_profile_dict_info['Node']['DataDir']:	
			bootup_options += ' --datadir=' + chain_profile_dict_info['Node']['DataDir']
		if chain_profile_dict_info['Plugin.MongoDB'] and chain_profile_dict_info['Plugin.MongoDB']['URL']:
			bootup_options += ' --mongodb=' + chain_profile_dict_info['Plugin.MongoDB']['URL']
		if chain_profile_dict_info['Delegate.SignKey']:
			pubkey = chain_profile_dict_info['Delegate.SignKey']['PublicKey']
			prikey = chain_profile_dict_info['Delegate.SignKey']['PrivateKey'] 
			bootup_options += ' --delegate-signkey=' + pubkey + ',' + prikey
		if chain_profile_dict_info['Rest'] and chain_profile_dict_info['Rest']['RESTHost'] and chain_profile_dict_info['Rest']['RESTPort']:
			bootup_options += ' --rest-servaddr=' + chain_profile_dict_info['Rest']['RESTHost']
			bootup_options += ' --restport=' + str(chain_profile_dict_info['Rest']['RESTPort'])
		if chain_profile_dict_info['Plugin.Wallet'] and chain_profile_dict_info['Plugin.Wallet']['WalletRESTHost'] and chain_profile_dict_info['Plugin.Wallet']['WalletRESTPort']:
			bootup_options += ' --wallet-rest-servaddr=' + chain_profile_dict_info['Plugin.Wallet']['WalletRESTHost']
			bootup_options += ' --wallet-rest-port=' + str(chain_profile_dict_info['Plugin.Wallet']['WalletRESTPort'])
		if chain_profile_dict_info['P2P'] and chain_profile_dict_info['P2P']['P2PServAddr'] and chain_profile_dict_info['P2P']['P2PPort'] and chain_profile_dict_info['P2P']['PeerList']:
			bootup_options += ' --p2p-servaddr=' + chain_profile_dict_info['P2P']['P2PServAddr']
			bootup_options += ' --p2pport=' + str(chain_profile_dict_info['P2P']['P2PPort'])
			bootup_options += ' --peerlist='
			
			peerlist = ''
			p2pport = chain_profile_dict_info['P2P']['P2PPort']
			for item in chain_profile_dict_info['P2P']['PeerList']:
				if peerlist:
					peer_url = ',' + item
					peerlist += peer_url
				
				else:
					peerlist = item + ':' 
			
			bootup_options += peerlist
	
	currpath = os.getcwd()
	bottos_startup_cmd = 'cd ' + BOTTOS_PROGRAM_WORK_DIR + '; ./bottos' + bootup_options + ' >core.file 2>&1'
	print bottos_startup_cmd
	common.do_cmd(bottos_startup_cmd, 'no_wait')
        common.do_cmd('cd ' + currpath)
	pass
    def node_build(self):
	pass	
    def node_stop(self, *args):
	if 'clean' in args:
		filename =  BOTTOS_PROGRAM_WORK_DIR+ '/datadir' 
		if os.path.exists(filename):
			shutil.rmtree(filename)
		
		filename = BOTTOS_PROGRAM_WORK_DIR+ '/core.file'
		if os.path.exists(filename):
			os.remove(filename)
	
	common.killbottos()
        pass




if __name__ == '__main__':
    _ = GOPATH
    _ = GLOBAL_BOTTOS_DIR
    _ = BOTTOS_PROGRAM_WORK_DIR		
    
	
    if len(sys.argv) <= 1 or sys.argv[1] in ['--help', '-h']:
	   common.print_help()
	   exit(1)
   
    conf = '.installation_config.txt'
    if os.path.exists(conf):
		with open(conf, 'r') as f:
			for line in f.readlines():
				if 'GLOBAL_BOTTOS_DIR' in line.split(':\n'):
					GLOBAL_BOTTOS_DIR = line.split(':\n')[1]
				elif 'BOTTOS_PROGRAM_WORK_DIR' in line.split(':\n'):
					BOTTOS_PROGRAM_WORK_DIR = line.split(':\n')[1]
    elif not sys.argv[1] == 'install':
	print '\n****** WARNING: you have no .installation_config.txt file under current directory, so use default GLOBAL_BOTTOS_DIR and BOTTOS_PROGRAM_WORK_DIR ***********\n'    

    if sys.argv[1] == 'install':
	   if os.geteuid() != 0:
                       print "The required packages must be installed under root user. Please turn into root account first."
                       exit(1)
	   x=raw_input('Do you need to install node by recommanded configuration? Y/N')
	   if x.upper() in ('Y', 'YES'):
		pass
	   elif x.upper() in ('N', 'NO'):
                
		x=raw_input('Please choose your working directory: default is %s ' %GLOBAL_BOTTOS_DIR )
		if len(x) > 0:   
			GLOBAL_BOTTOS_DIR = x
		x=raw_input('Please choose your bottos programe files\' directory: default is %s ' % BOTTOS_PROGRAM_WORK_DIR )   
		if len(x) > 0:
			BOTTOS_PROGRAM_WORK_DIR = x
		   
		if not os.path.exists(BOTTOS_PROGRAM_WORK_DIR):
			os.makedirs(BOTTOS_PROGRAM_WORK_DIR)
		   
		x=raw_input('Do you need to install base packages? Y/N: default is %s ' % user_choice_list['install_base']) 
		if len(x) > 0:
			if x.upper() in ('Y', 'YES'):   
				user_choice_list['install_base'] = 'yes'
			elif x.upper() in ('N', 'NO'):
				user_choice_list['install_base'] = 'no'
		else:
			pass
		   
		x=raw_input('Do you need to install golang? Y/N: default is %s ' % user_choice_list['install_golang']) 
		if len(x) > 0:
			if x.upper() in ('Y', 'YES'):   
				user_choice_list['install_golang'] = 'yes'
			elif x.upper() in ('N', 'NO'):
				user_choice_list['install_golang'] = 'no'
		else:
			pass
		   
		x=raw_input('Do you need to install mongodb? Y/N: default is %s ' % user_choice_list['install_mongodb']) 
		if len(x) > 0:
			if x.upper() in ('Y', 'YES'):   
				user_choice_list['install_mongodb'] = 'yes'
			elif x.upper() in ('N', 'NO'):
				user_choice_list['install_mongodb'] = 'no'
		else:
			pass
		   
		x=raw_input('Do you need to install bottos source code? Y/N: default is %s ' % user_choice_list['install_bottos_source_code']) 
		if len(x) > 0:
			if x.upper() in ('Y', 'YES'):   
				user_choice_list['install_bottos_source_code'] = 'yes'
			elif x.upper() in ('N', 'NO'):
				user_choice_list['install_bottos_source_code'] = 'no'
		else:
			pass
	   else:
		print 'Wrong input. Please try again.'
		exit(1)
	   		   
	   	
	   print 'Your install configuration is as following: \n'	   	
	   for key, item in user_choice_list.items():
		if item is 'yes':
			print key
			print
	   	
	   print 'Your working directory: %s' %GLOBAL_BOTTOS_DIR
           print 'Your bottos programe files\' directory: %s' % BOTTOS_PROGRAM_WORK_DIR
	   print
	   
           x=raw_input('Are you sure? Y/N')
	   if  x.upper() in ('Y', 'YES'):
		
	    	print '\n============ NOW INSTALL BOTTOS ENVIRONMENT =======================\n'
	   	deploy_obj = bottos_node_deploy()
	   	deploy_obj.install_env()
	   else:
		print 'All right. Bye bye.'
		exit(1)	
   
    elif sys.argv[1] == 'build':
	    print '\n============ NOW BUILD BOTTOS =======================\n'
	    node_build_obj = bottos_node_build()
	    #node_build_obj.build_bottos()
	    node_build_obj.build_bottos_release()
		
    elif sys.argv[1] == 'start':
	    
	    #print '\n============ NOW DOWNLOAD BOTTOS =======================\n'
	    #deploy_obj = bottos_node_deploy()
	    #deploy_obj.install_env()
	    
	    #exit(1)
	    #print '\n============ NOW BUILD BOTTOS =======================\n'
	    #node_build_obj = bottos_node_build()
	    #node_build_obj.build_bottos()

	    node_profile_obj = bottos_node_profile()
	    node_profile_obj.generate_default_profiles()
	    
	    node_profile_obj.show_profiles()
             
	    #node_profile_obj.show_profiles( 'node_profile_info.toml', 'service_profile_info.toml')    
	    
	    node_profile_obj.set_profile_info('chain_profile_info.toml')
	    	
	    print '\n============ NOW CHECK PROFILE OF MONGO DB =======================\n'
		
	    node_profile_obj.set_profile_info('mongodb_profile_info.toml')
	    
	    print '\n============ NOW CHECK PROFILE OF BOOT UP =======================\n'
		 
	    node_profile_obj.set_profile_info('bottos_bootup_options_profile_info.toml')
		
	    print '\n============ NOW START BOTTOS =======================\n'
	    node_profile_obj.show_profiles()
	    
	    #common.download_official_bcli()
            
	    if node_profile_obj.mode is 'to_bottos_net':
		common.download_official_genesis()
	    
	    node_start = bottos_node_apply()
	    node_start.node_start('clean')
	    	
    elif sys.argv[1] == 'stop':
	    node_start = bottos_node_apply()
	    node_start.node_stop('clean')
    elif sys.argv[1] == 'show':
	    node_profile_obj = bottos_node_profile()
	    node_profile_obj.show_profiles()	
    else:
    	#obj = download_progress_bar('https://studygolang.com/dl/golang/go1.10.1.linux-amd64.tar.gz') 
    	#obj.download_with_progressbar('./go1.10.1.linux-amd64.tar.gz')
	common.print_help()

    exit(0)


