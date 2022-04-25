package basefunc

func PageLimit(count, page, limit int) (begin, end int) {
	if page <= 0 || limit <= 0 || count <= 0 {
		return
	}
	if (page-1)*limit >= count {
		return
	}
	if page*limit >= count {
		end = count
	} else {
		end = page * limit
	}
	if page != 0 {
		begin = (page - 1) * limit
	}
	return
}
