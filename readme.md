## Solution:

### Summary
To make the update be atomically, we have to lock the element `i, i+1, i+2 and j`
since only `j` is the updated one, so we can lock `j` with write-lock, `i, i+1, i+2` with read-lock.

Each can do update only if it holds the locks of the 4 elements, and when the update is finished, work should release the locks.

### Design

Define an `lock` array that has the same length with `S`, and the type is `int32`. `lock[i]` stand for the lock status for the
element `i`, the value can be: 
1. `0`: no locks on this element.
2. A big number(`2<<20`): write-lock on this element.
3. other value: read-lock on this element.

On lock phase, update `lock[i] = 2<<20` to add a write-lock iif `lock[i]` is `0`, 
update `lock[i] = lock[i]+1` to add a read-lock if `lock[i]` is not `2<<20`.

On unlock phase, update `lock[i] = 0` to unlock a write-lock, update `lock[i] = lock[i] - 1` to unlock
a read-lock.

To prevent the deadlock, we can sort the `i, i+1, i+2, j` and lock the elements from small to big, if there is someone that
can not be locked, just wait the lock owner release it, there will never deadlock.

### Test

To test the system, we set a _shadow_ `S_s`, it is a copy from the original `S`, for each worker, when doing the update, we
record the `i, j` and reply the update on `S_s`, at the end, `S_s` and `S` should be the same.