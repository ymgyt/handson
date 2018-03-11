#
# Cookbook:: mysql
# Recipe:: default
#
# Copyright:: 2018, The Authors, All Rights Reserved.

package "mysql-server" do
  action :install
end

service "mysqld" do
  action [ :enable, :start ]
end
