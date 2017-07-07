--[[entry
	psharingrate 代理分成占比
	price 当前价格
]]-- 

function entry(count)
	local r = tonumber(count)
	if(r == nil)
	then
		return 0
	else
		return r
	end
end

--[[ test
print (entry(0,10,10,10,10000))
print (entry(1,10,10,10,10000))
]]--