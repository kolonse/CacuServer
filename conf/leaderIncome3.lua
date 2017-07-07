--[[entry
	计算会长应得的总收益,
	如果是推广人是会长,那么按照会长分成算
	如果是副会长,那么需要计算出会长应该得到的分成
]]-- 

function entry(utype, usharingrate,psharingrate,ursharingrate,price)
	utype = tonumber(utype)
	usharingrate = tonumber(usharingrate)
	psharingrate = tonumber(psharingrate)
	ursharingrate = tonumber(ursharingrate)
	price = tonumber(price)
	if(utype == 0)
	then
		return price * psharingrate * usharingrate / 10000
	else
		return price * (psharingrate / 100) * (1 - usharingrate/100) * (ursharingrate/100)
	end
end

--[[ test
print (entry(0,10,10,10,10000))
print (entry(1,10,10,10,10000))
]]--