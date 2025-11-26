-- Shouldn't work, cause of recursion.
fn pow(x int, n int) int then
    if n == 0 then
        return 1
    end

    return x * pow(x, n - 1)
end

let result int = pow(2, 8)

print(stri(result))

