local name, website, description = unpack(ARGV)
local companyID = redis.call('zrange', 'companies', -1, -1)[1]

--generate new company id
if companyID == nil then companyID = 0 end
companyID = companyID + 1

--enter the company into the main list
redis.call('zadd', 'companies', companyID, companyID)

--store their data
redis.call('hmset', 'companies:' .. companyID, 'name', name, 'website', website, 'description', description)

return companyID