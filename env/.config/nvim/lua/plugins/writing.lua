return {
	{
		"epwalsh/obsidian.nvim",
		version = "*",
		lazy = true,
		ft = "markdown",
		config = function()
			require("obsidian").setup({
				workspaces = {
					{
						name = "vault",
						path = "~/vault",
					},
				},
				notes_subdir = ".",
				new_notes_location = ".",
				new_note_filename = function(title)
					return title:gsub(" ", "_") -- Replace spaces with underscores
				end,

				disable_frontmatter = false,
				-- key mappings, below are the defaults
				mappings = {
					-- overrides the 'gf' mapping to work on markdown/wiki links within your vault
					["gd"] = {
						action = function()
							return require("obsidian").util.gf_passthrough()
						end,
						opts = { noremap = false, expr = true, buffer = true },
					},
					-- toggle check-boxes
					["<leader>ti"] = {
						action = function()
							return require("obsidian").util.toggle_checkbox()
						end,
						opts = { buffer = true },
					},
				},
				completion = {
					nvim_cmp = true,
					min_chars = 2,
				},
				ui = {
					-- Disable some things below here because I set these manually for all Markdown files using treesitter
					checkboxes = {},
					bullets = {},
				},
			})
		end,
		dependencies = {
			"nvim-lua/plenary.nvim",
		},
	},
}
