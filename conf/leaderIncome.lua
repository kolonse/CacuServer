--[[entry
	utype 推荐人类型 会长/副会长
	usharingrate 会长/副会长 分成占比
	psharingrate 代理分成占比
	ursharingrate 工会会长分成占比
	price 当前价格
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
		return price * (psharingrate / 100) * (usharingrate/100) * (ursharingrate/100)
	end
end

--[[ test
print (entry(0,10,10,10,10000))
print (entry(1,10,10,10,10000))
]]--