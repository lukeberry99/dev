-- Custom highlight groups for consistent theming
local function set_highlights()
	local highlights = {
		-- Blink.cmp highlights
		BlinkCmpMenu = { bg = "NONE", fg = "NONE" },
		BlinkCmpMenuBorder = { bg = "NONE", fg = "#6c7086" },
		BlinkCmpMenuSelection = { bg = "#313244", fg = "NONE" },
		BlinkCmpDoc = { bg = "NONE", fg = "NONE" },
		BlinkCmpDocBorder = { bg = "NONE", fg = "#6c7086" },
		BlinkCmpDocCursorLine = { bg = "#313244", fg = "NONE" },
		BlinkCmpSignatureHelp = { bg = "NONE", fg = "NONE" },
		BlinkCmpSignatureHelpBorder = { bg = "NONE", fg = "#6c7086" },

		-- Snacks.nvim picker highlights
		SnacksPickerBorder = { bg = "NONE", fg = "#6c7086" },
		SnacksPickerTitle = { bg = "NONE", fg = "#f38ba8" },
		SnacksPickerHeader = { bg = "NONE", fg = "#94e2d5" },
		SnacksPickerCursor = { bg = "#313244", fg = "NONE" },
		SnacksPickerSelected = { bg = "#313244", fg = "NONE" },

		-- Noice highlights
		NoiceCmdline = { bg = "NONE", fg = "NONE" },
		NoiceCmdlineIcon = { bg = "NONE", fg = "#89b4fa" },
		NoiceCmdlinePopupBorder = { bg = "NONE", fg = "#6c7086" },
		NoiceCmdlinePopupTitle = { bg = "NONE", fg = "#f38ba8" },

		-- Diagnostic highlights
		DiagnosticError = { fg = "#f38ba8" },
		DiagnosticWarn = { fg = "#fab387" },
		DiagnosticInfo = { fg = "#89b4fa" },
		DiagnosticHint = { fg = "#94e2d5" },
		DiagnosticUnderlineError = { sp = "#f38ba8", undercurl = true },
		DiagnosticUnderlineWarn = { sp = "#fab387", undercurl = true },
		DiagnosticUnderlineInfo = { sp = "#89b4fa", undercurl = true },
		DiagnosticUnderlineHint = { sp = "#94e2d5", undercurl = true },

		-- Which-key highlights
		WhichKeyBorder = { bg = "NONE", fg = "#6c7086" },
		WhichKeyTitle = { bg = "NONE", fg = "#f38ba8" },
		WhichKeyGroup = { bg = "NONE", fg = "#94e2d5" },
		WhichKeyDesc = { bg = "NONE", fg = "#cdd6f4" },
		WhichKeySeperator = { bg = "NONE", fg = "#6c7086" },
		WhichKeyFloating = { bg = "NONE", fg = "NONE" },

		-- General UI highlights
		FloatBorder = { bg = "NONE", fg = "#6c7086" },
		FloatTitle = { bg = "NONE", fg = "#f38ba8" },
		NormalFloat = { bg = "NONE", fg = "NONE" },
		Pmenu = { bg = "NONE", fg = "NONE" },
		PmenuSel = { bg = "#313244", fg = "NONE" },
		PmenuBorder = { bg = "NONE", fg = "#6c7086" },
	}

	for group, opts in pairs(highlights) do
		vim.api.nvim_set_hl(0, group, opts)
	end
end

-- Apply highlights after colorscheme is loaded
vim.api.nvim_create_autocmd("ColorScheme", {
	callback = set_highlights,
})

-- Apply highlights immediately if already loaded
if vim.g.colors_name then
	set_highlights()
end