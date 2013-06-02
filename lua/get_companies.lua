local companiesList = redis.call('zrange', 'companies', 0, -1)
local companiesData = {}

for companyIdx, companyID in pairs(companiesList) do
	local companyData = redis.call('hgetall', 'companies:'..companyID)
	table.insert(companiesData, companyData)
end

return companiesData