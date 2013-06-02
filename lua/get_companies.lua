local companiesList = redis.call('zrange', 'companies', 0, -1)
local companiesData = {}

for companyIdx, companyID in pairs(companiesList) do
	local companyData = redis.call('hgetall', 'companies:'..companyID)

	for dataIdx, data in pairs(companyData) do
		redis.log(redis.LOG_WARNING, dataIdx, dataIdx % 2)

		if dataIdx % 2 == 0 then
			table.insert(companiesData, data)
		end
	end
end

return companiesData