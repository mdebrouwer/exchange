# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = '2'

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = 'ubuntu/yakkety64'
  config.vm.synced_folder '.', '/home/ubuntu/workspace/src/github.com/mdebrouwer/exchange'
  config.ssh.forward_agent = true
  config.vm.provision 'file', source: '~/.gitconfig', destination: '.gitconfig'

  config.vm.provision 'shell', inline: 'curl -sL https://deb.nodesource.com/setup_7.x | bash -'
  config.vm.provision 'shell', inline: 'apt-get --yes update'
  config.vm.provision 'shell', inline: 'apt-get --yes upgrade'

  [
    'docker.io',
    'nodejs',
    'make',
    'vim'
  ].each do |pkg|
    config.vm.provision 'shell', inline: "apt-get --yes install #{pkg}"
  end

  config.vm.provision 'shell', inline: 'usermod -a -G docker ubuntu'

  config.vm.provider :virtualbox do |vb|
    vb.customize ['modifyvm', :id, '--memory', '4096']
  end
end
