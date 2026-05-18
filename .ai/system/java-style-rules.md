# JAVA STYLE RULES

==================================================
IMPORT RULES
==================================================

NEVER use fully-qualified class names inline in method signatures,
field declarations, generics, lambdas, or implementations.

BAD:
(java.util.concurrent.Callable<T> callable)

BAD:
private java.util.Map<String, Object> data;

BAD:
new java.util.ArrayList<>();

GOOD:
import java.util.concurrent.Callable;
import java.util.Map;
import java.util.ArrayList;

Use:
Callable<T>
Map<String, Object>
ArrayList<>()

==================================================
MANDATORY IMPORT ORGANIZATION
==================================================

ALWAYS:
- generate proper imports
- organize imports cleanly
- remove unused imports
- use IDE-quality import formatting

Imports MUST:
- be explicit
- be readable
- follow production Java conventions

==================================================
EXCEPTIONS
==================================================

Fully-qualified names are ONLY allowed:
- when resolving naming conflicts
- in generated code
- in reflection-heavy infrastructure code
- when explicitly required

Otherwise:
ALWAYS import classes properly.

==================================================
CODE QUALITY ENFORCEMENT
==================================================

Generated Java code MUST resemble:
- senior production engineer code
- clean enterprise code
- maintainable codebases

NOT:
- compiler-generated style
- IDE fallback style
- decompiled style
- generated boilerplate style