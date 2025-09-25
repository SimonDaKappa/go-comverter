This file contains examples of comments found throughout , and the comment formats that define
each case

# Doxygen Comments
See: [Doxygen Comment SPEC](https://www.doxygen.nl/manual/docblocks.html)

# Header/Footer Formats

## C/C++

### Javadoc
You can use the Javadoc style, which consist of a C-style comment block starting with 2+ *'s, like this:
```md
/**
 * ... text ...
 */
```

or 
```md
/**************************************************************************
 * ... text ...
 */
```

or 
```md
/**************************************************************************
 * ... text ...
 **************************************************************************/
```

**Example:**
```cpp
/**
 * @name foo
 * @brief foo(l)'s around
 * @param[in] bar
 * @param[out] guh
 * 
 * @return 0 on success, -1 on failure.
 */
```

### QT
You can use the Qt style and add an exclamation mark (!) after the opening of a C-style comment block, as shown in this example:
```md
/*!
 * ... text ...
 */
```

### Triple Slash/Excl.
A block of at least two C++ comment lines, where each line starts with an additional slash or an exclamation mark. 

Here are examples of the two cases:

```md
///
/// ... text ...
///
```

or 

```md
//////////////////////////////////////////////////////////////////////////
/// ... text ...
//////////////////////////////////////////////////////////////////////////
```

or

```md
//!
//!... text ...
//!
```
### Block Comment Visibility Overrides
Override custom comment header/footer styles (used for visibility) with doxygen tags
```md
                                            ↓ This is the override.
/*******************************************//**
 *  ... text
 ***********************************************/ <- The single line of "*" repeated is treated as empty lines, until we hit "*/" ?
```

## Body Formats

### 

