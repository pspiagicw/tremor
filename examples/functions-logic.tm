-- currently can't determine embedded return statements.
fn max(a int, b int) int then
    if a > b then
        return a
    end
end

print(stri(max(10, 3)))

