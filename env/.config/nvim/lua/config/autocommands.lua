vim.api.nvim_create_autocmd("TextYankPost", {
	desc = "Highlight when yanking text",
	group = vim.api.nvim_create_augroup("highlight-yank", { clear = true }),
	callback = function()
		vim.highlight.on_yank()
	end,
})

vim.api.nvim_create_autocmd("FileType", {
	pattern = "markdown",
	callback = function()
		vim.opt_local.wrap = true -- enables visual line wrapping
		vim.opt_local.linebreak = true -- wraps at word boundaries
		vim.opt_local.colorcolumn = "80" -- visual guide at 80 chars
		vim.opt_local.textwidth = 80 -- hard wrap via formatting
	end,
})
