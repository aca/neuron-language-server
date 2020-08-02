# neuron-language-server
Language server for [neuron](https://github.com/srid/neuron).

Neuron will embed language server in neuron. Check [sric/neuron#213](https://github.com/srid/neuron/issues/213) for updates.<br/> 
This is just an personal experiment.

Suports
- textDocument/completion(search by title/id)<br/>
  ![completion](./images/completion.png)
- textDocument/definition
- textDocument/hover<br/>
  ![definition](./images/definition.png)


#### TODO
- improve completion
- textDocument/codeAction
- .... a lot, I don't know when it will be


#### Prerequisites
  - neuron

#### Installation
```
go get -u github.com/aca/neuron-language-server
# Or: nix-build (if you have Nix)
```

#### LSP client settings
- vim/neovim, [coc.nvim](https://github.com/neoclide/coc.nvim)
  ```
  "languageserver": {
    "neuron": {
      "command": "neuron-language-server",
      "filetypes": ["markdown"]
    },
  ```
- neovim, [nvim-lsp](https://github.com/neovim/nvim-lsp)
  ```lua
  local nvim_lsp = require('nvim_lsp')
  local configs = require('nvim_lsp/configs')

  configs.neuron_ls = {
    default_config = {
      cmd = {'neuron-language-server'};
      filetypes = {'markdown'};
      root_dir = function()
        return vim.loop.cwd()
      end;
      settings = {};
    };
  }
  nvim_lsp.neuron_ls.setup{}
  ```
