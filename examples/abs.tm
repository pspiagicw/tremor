fn abs(x int) int then
    if x < 0 then
        return 0 - x
    else
        return x
    end
end

-- prefix expression doesn't work I think
let result int = abs(-42)

