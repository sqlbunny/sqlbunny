package gen

import "github.com/KernelPay/sqlboiler/schema"

func buildImports(fields []*schema.Field) []string {
	var res []string
	for _, f := range fields {
		if t, ok := f.Type.(schema.TypeWithImports); ok {
			res = append(res, t.GetImports()...)
		}
	}
	return removeDuplicates(res)
}

func removeDuplicates(dedup []string) []string {
	if len(dedup) <= 1 {
		return dedup
	}

	for i := 0; i < len(dedup)-1; i++ {
		for j := i + 1; j < len(dedup); j++ {
			if dedup[i] != dedup[j] {
				continue
			}

			if j != len(dedup)-1 {
				dedup[j] = dedup[len(dedup)-1]
				j--
			}
			dedup = dedup[:len(dedup)-1]
		}
	}

	return dedup
}
