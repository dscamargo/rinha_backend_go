package pessoasdb

const (
	QuerySelectPessoaById    = "SELECT id, nome, apelido, nascimento, stack FROM public.pessoas WHERE id = $1;"
	QuerySelectPessoasByTerm = "SELECT id, nome, apelido, nascimento, stack FROM public.pessoas p WHERE search ILIKE '%' || $1 || '%' LIMIT 50;"
	QueryCountAllPessoas     = "SELECT count(*) FROM public.pessoas;"
)
