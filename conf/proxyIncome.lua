--[[entry
	psharingrate 代理分成占比
	ursharingrate 工会会长分成占比
	price 当前价格
]]-- 

function entry(psharingrate,ursharingrate,price)
	psharingrate = tonumber(psharingrate)
	ursharingrate = tonumber(ursharingrate)
	price = tonumber(price)
	return price*(psharingrate/100)*(100-ursharingrate)/100
end

--[[ test
print (entry(0,10,10,10,10000))
print (entry(1,10,10,10,10000))
]]--