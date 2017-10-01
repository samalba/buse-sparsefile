# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "debian/stretch64"

  config.vm.network "private_network", type: "dhcp"
  config.vm.network "public_network"

  config.vm.synced_folder ".", "/home/vagrant/go/buse-sparsefile"

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "1024"
  end

  # The following two lines require the vbguest plugin to be installed:
  # vagrant plugin install vagrant-vbguest
  config.vbguest.auto_update = true
  config.vbguest.auto_reboot = true

  config.vm.provision "shell", path: "scripts/configure_vm.sh", env: {"VAGRANT" => "true"}
end
