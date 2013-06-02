local companiesList = redis.call('zrange', 'companies', 0, -1)
local companiesData = {}

for idx, companyID in pairs(companiesList) do
	local companyData = redis.call('hgetall', 'companies:'..companyID)
	companiesData[idx] = companyData
end

return companiesData