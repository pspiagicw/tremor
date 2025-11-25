-- currently can't determine embedded return statements.
fn max(a int, b int) int then
    if a > b then
        return a
    end
    return b
end

print(stri(max(10, 3)))

