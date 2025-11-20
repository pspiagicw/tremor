-- currently can't determine embedded return statements.
fn max(a int, b int) int then
    if a > b then
        return a
    else
        return b
    end
    return 1
end

print(str(max(10, 3)))

