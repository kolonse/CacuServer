--[[entry
	psharingrate 代理分成占比
	price 当前价格
]]-- 

function entry(psharingrate,price)
	psharingrate = tonumber(psharingrate)
	price = tonumber(price)
	return price*(100-psharingrate)/100
end

--[[ test
print (entry(0,10,10,10,10000))
print (entry(1,10,10,10,10000))
]]--