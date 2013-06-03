local companiesList = redis.call('zrange', 'companies', 0, -1)
local companiesData = {}

for companyIdx, companyID in pairs(companiesList) do
	table.insert(companiesData, {redis.call('hgetall', 'companies:'..companyID), redis.call('smembers', 'companies:'..companyID..':languages')})
end

return companiesData