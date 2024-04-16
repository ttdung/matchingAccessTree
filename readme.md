## problem
Define at [OpenABE API Doc](https://github.com/zeutro/openabe/blob/master/docs/libopenabe-v1.0.0-api-doc.pdf) - Section 2.3 - page 13.
### Policy Tree Structure
- Policy tree is a (**boolean formulas**). Include:   
  1. OR
  2. AND
  3. Comparations: `<, >, <=, >=, =`,   
    a.  Date Comparison: days since the beginning of unix epoch time (Jan 1, 1970).  
       - Example:  `Date > March 1, 2015`
  4. DATE: `[Prefix] = [Month] [Day], [Year]`
  5. RANGE:   
    a. Integer range: `[Attribute] in ([int] - [int])`.    (Attribute > int1 and Attribute < int2)
    b. Date Range:    `[Attribute] = [Month] [Day] - [Day], [Year]`  (Attribute >= date1 and Attribute <= date2)
  6. Brake `()`

### Attribute Lists
- Separated by `|`. 
- Example: `|IT|Manager|Experience=5|Date = December 20, 2015|`

### Requirement
- Input:   
  1. `Attribute Lists` String. 
  2. `Policy Tree` String.

- Output: True/ False

## Algorithm

### Using stack and boolean formulas 
- Use 2 Stack: 
  1. Stack contains operators: `operators`
  2. Stack contains true/false of the statement: `values`

- Read an element in `Policy Tree`:
  1. IF `(`: Push into `operators`
  2. IF `)`: Pop all `operators` and `values` to execute operators until reach `(`. 
      - Push result back to `values`
      - Pop `(` from `operators`
  3. IF `AND/and`, `OR/or`: Push into `operators`
  4. IF `boolean formulas`: solve boolean fomulas and push into `values`.

