--[[entry

]]-- 

function entry(curprice,totalprice)
	curprice = tonumber(curprice)
	totalprice = tonumber(totalprice)
	if totalprice == nil or totalprice == false
	then
		totalprice = 0
	end
	if curprice == nil or curprice == false
	then
		return 0
	end
	return curprice - totalprice
end

--[[ test
print (entry(0,10,10,10,10000))
print (entry(1,10,10,10,10000))
print (entry(10,xx))
]]--
