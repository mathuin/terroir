# -*- mode: ruby -*-
# vi: set ft=ruby :

# based heavily on nathany/vagrant-gopher
def src_path
  File.dirname(__FILE__)
end

def bootstrap()
  go_archive = "go1.3.1.linux-amd64.tar.gz"
  gdal_version = "1.11.0"
  gdal_dir = "gdal-#{gdal_version}"
  gdal_archive = "#{gdal_dir}.tar.gz"

  profile = <<-PROFILE
    export GOPATH=$HOME
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
    export CDPATH=.:$GOPATH/src/github.com:$GOPATH/src/code.google.com/p:$GOPATH/src/bitbucket.org:$GOPATH/src/launchpad.net
  PROFILE

  script = <<-SCRIPT
    apt-get -qq update
    apt-get -qq upgrade
    apt-get -qq install git mercurial bzr curl
    if ! [ -f /home/vagrant/#{go_archive} ]; then
      response=$(curl -Os https://storage.googleapis.com/golang/#{go_archive})
    fi
    tar -C /usr/local -xzf #{go_archive}

    apt-get -qq install build-essential python-all-dev pkg-config proj-bin libproj-dev
    if ! [ -f /home/vagrant/#{gdal_archive} ]; then
      response=$(curl -Os http://download.osgeo.org/gdal/#{gdal_version}/#{gdal_archive})
    fi
    tar zxf #{gdal_archive}
    cd ./#{gdal_dir}
    ./configure --prefix=/usr --with-python
    make -j8
    make install

    echo '#{profile}' >> /home/vagrant/.profile

    date > /etc/vagrant_provisioned_at
  SCRIPT
end

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.hostname = "golang"
  config.vm.synced_folder "../../..", "/home/vagrant/src"
  config.vm.provision :shell, :inline => bootstrap()

  config.vm.provider "virtualbox" do |v|
    v.memory = 1024
    v.cpus = 4
  end
end
