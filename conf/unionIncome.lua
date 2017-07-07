--[[entry
	utype 推荐人类型 会长/副会长
	usharingrate 会长/副会长 分成占比
	psharingrate 代理分成占比
	ursharingrate 工会会长分成占比
	price 当前价格
]]-- 

function entry(psharingrate,ursharingrate,price)
	psharingrate = tonumber(psharingrate)
	price = tonumber(price)
	return price*ursharingrate*psharingrate / 10000
end

--[[ test
print (entry(0,10,10,10,10000))
print (entry(1,10,10,10,10000))
]]--