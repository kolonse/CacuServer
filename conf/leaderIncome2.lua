--[[entry
	计算 推广人分成给会长的部分
	如果是会长,那么应该是 0
]]-- 

function entry(utype, usharingrate,psharingrate,ursharingrate,price)
	utype = tonumber(utype)
	usharingrate = tonumber(usharingrate)
	psharingrate = tonumber(psharingrate)
	ursharingrate = tonumber(ursharingrate)
	price = tonumber(price)
	if(utype == 0)
	then
		return 0
	else
		return price * (psharingrate / 100) * (1 - usharingrate/100) * (ursharingrate/100)
	end
end

--[[ test
print (entry(0,10,10,10,10000))
print (entry(1,10,10,10,10000))
]]--