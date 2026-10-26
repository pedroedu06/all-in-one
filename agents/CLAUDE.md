AGENTS RULES

Development in stages, with a plan in place before any code is written.

1 - Think before you implement Declare → Confirm → Implement

2 - Never use any in typescript

// Forbidden function process(data: any) { ... }

// Correct function process(data: unknown) { ... } function process(data: T) { ... }

3 - Simplicity first Write the minimum amount of code that correctly solves the problem.

3 — Surgical Changes Touch only what is necessary to accomplish the task.

6 — Read Before You Write Understand the existing code before modifying it.

7 — Token Budgets Are Not Advisory Token budgets and constraints are hard limits, not suggestions.

8 — Tests Must Verify Intent, Not Just Behavior Tests must validate why the code exists, not only what it does.

9 — Fail Loud Errors must be visible, explicit, and informative.

10 - never read .env files, only test and check the log. if erro is in credencials or something else return a message for me to correct the .env - dont read it 
