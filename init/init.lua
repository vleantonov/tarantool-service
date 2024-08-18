local username = os.getenv('TARANTOOL_USER_NAME')
local password = os.getenv('TARANTOOL_USER_PASSWORD')
local port = os.getenv('TARANTOOL_PORT')

box.cfg{ listen = port }

box.schema.user.create(username, { password = password, if_not_exists = true })
box.once('grant_user', function()
    box.schema.user.grant(username, 'read,write,execute', 'universe')
end)

-- Create user space --
local users = box.schema.create_space('users', { format = {
    { name = 'username', type = 'string' },
    { name = 'password', type = 'string' }
}, if_not_exists = true })


users:create_index('pk', { parts = {'username'}, if_not_exists = true })
users:insert{ 'admin', 'presale' }

-- Create data space --
local data = box.schema.create_space('data', { format = {
    { name = 'key', type = 'string' },
    { name = 'value' }
}, if_not_exists = true })

data:create_index('pk', { parts = { 'key' }, if_not_exists = true, type = 'HASH' })
