# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = '2'

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = 'ubuntu/yakkety64'
  config.vm.network 'forwarded_port', guest: 6288, host: 6288
  config.vm.network 'forwarded_port', guest: 8080, host: 8080

  config.vm.synced_folder '.', '/home/ubuntu/workspace/src/github.com/mdebrouwer/exchange'
  config.vm.provision 'shell', inline: 'chown -R ubuntu:ubuntu /home/ubuntu/workspace'

  config.ssh.forward_agent = true
  config.vm.provision 'file', source: '~/.gitconfig', destination: '.gitconfig'

  config.vm.provision 'shell', inline: 'curl -sL https://deb.nodesource.com/setup_7.x | bash -'
  config.vm.provision 'shell', inline: 'apt-get --yes update'
  config.vm.provision 'shell', inline: 'apt-get --yes upgrade'

  [
    'docker.io',
    'nodejs',
    'golang-1.7',
    'make',
    'vim'
  ].each do |pkg|
    config.vm.provision 'shell', inline: "apt-get --yes install #{pkg}"
  end

  config.vm.provision 'shell', inline: 'usermod -a -G docker ubuntu'
  config.vm.provision 'shell', inline: 'ln -s /usr/lib/go-1.7/bin/go /bin/go'
  config.vm.provision 'shell', inline: 'echo "export GOPATH=/home/ubuntu/workspace" >> /home/ubuntu/.profile'

  config.vm.provider :virtualbox do |vb|
    vb.customize ['modifyvm', :id, '--memory', '4096']
  end
end
