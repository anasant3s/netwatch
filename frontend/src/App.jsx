
import { useRef, useState } from 'react'
import {
  Activity,
  Bot,
  ChevronDown,
  CircleHelp,
  FileText,
  FolderOpen,
  Menu,
  MessageSquare,
  MoreHorizontal,
  Paperclip,
  Plus,
  Search,
  Send,
  Settings,
  ShieldCheck,
  Sparkles,
  Upload,
  X,
} from 'lucide-react'

const initialConversations = [
  { id: 1, title: 'Especificações do NVR5464-EI' },
  { id: 2, title: 'Configuração de câmeras IP' },
  { id: 3, title: 'Diagnóstico de conectividade' },
]

const suggestions = [
  {
    icon: FileText,
    title: 'Consultar especificações',
    description: 'Encontre informações técnicas nos seus documentos.',
    prompt: 'Quais são as principais especificações do NVR5464-EI?',
  },
  {
    icon: Activity,
    title: 'Analisar equipamentos',
    description: 'Entenda características e funcionalidades dos dispositivos.',
    prompt: 'Quais funcionalidades de análise inteligente estão disponíveis?',
  },
  {
    icon: ShieldCheck,
    title: 'Tirar dúvidas técnicas',
    description: 'Consulte manuais e documentações técnicas.',
    prompt: 'Como funciona a configuração de rede do equipamento?',
  },
]

function App() {
  const [conversations, setConversations] = useState(initialConversations)
  const [activeConversation, setActiveConversation] = useState(null)
  const [messages, setMessages] = useState([])
  const [question, setQuestion] = useState('')
  const [selectedFile, setSelectedFile] = useState(null)
  const [notice, setNotice] = useState('')
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const fileInputRef = useRef(null)

  function startNewConversation() {
    setActiveConversation(null)
    setMessages([])
    setQuestion('')
    setNotice('')
  }

  function selectConversation(conversation) {
    setActiveConversation(conversation)
    setMessages([])
    setNotice('A integração com o histórico do backend será adicionada na próxima etapa.')
  }

  function handleFileChange(event) {
    const file = event.target.files?.[0]

    if (!file) return

    if (!file.name.toLowerCase().endsWith('.pdf')) {
      setNotice('Selecione um arquivo PDF.')
      return
    }

    setSelectedFile(file)
    setNotice(`Arquivo selecionado: ${file.name}. O envio será conectado ao backend em seguida.`)
  }

  function handleSubmit(event) {
    event.preventDefault()

    const trimmedQuestion = question.trim()

    if (!trimmedQuestion) return

    setMessages((currentMessages) => [
      ...currentMessages,
      {
        id: Date.now(),
        role: 'user',
        content: trimmedQuestion,
      },
    ])

    if (!activeConversation) {
      const newConversation = {
        id: Date.now(),
        title: trimmedQuestion.length > 38
          ? `${trimmedQuestion.slice(0, 38)}...`
          : trimmedQuestion,
      }

      setActiveConversation(newConversation)
      setConversations((current) => [newConversation, ...current])
    }

    setQuestion('')
    setNotice('Mensagem registrada na interface. A conexão com o Gemini será implementada na próxima etapa.')
  }

  return (
    <div className="flex h-screen overflow-hidden bg-[#090d16] text-slate-100">
      {sidebarOpen && (
        <aside className="flex w-[280px] shrink-0 flex-col border-r border-white/[0.07] bg-[#0d121e] max-md:absolute max-md:inset-y-0 max-md:left-0 max-md:z-30 max-md:w-[290px]">
          <div className="flex h-[76px] items-center justify-between border-b border-white/[0.07] px-5">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-cyan-400/10 text-cyan-400">
                <Activity size={22} strokeWidth={2.2} />
              </div>
              <div>
                <h1 className="text-[16px] font-semibold tracking-wide">
                  NetWatch <span className="text-cyan-400">AI</span>
                </h1>
                <p className="mt-0.5 text-[10px] tracking-[0.18em] text-slate-500">
                  INTELLIGENT NETWORK
                </p>
              </div>
            </div>

            <button
              onClick={() => setSidebarOpen(false)}
              className="rounded-lg p-2 text-slate-500 transition hover:bg-white/5 hover:text-white"
              aria-label="Fechar menu"
            >
              <X size={17} />
            </button>
          </div>

          <div className="p-4">
            <button
              onClick={startNewConversation}
              className="flex w-full items-center justify-center gap-2 rounded-xl bg-cyan-400 px-4 py-3 text-sm font-semibold text-slate-950 transition hover:bg-cyan-300"
            >
              <Plus size={17} />
              Nova conversa
            </button>

            <div className="mt-5 flex items-center gap-2 rounded-lg border border-white/[0.07] bg-white/[0.025] px-3 py-2.5 text-slate-500">
              <Search size={16} />
              <span className="text-xs">Buscar conversas</span>
              <span className="ml-auto text-[10px] text-slate-600">⌘ K</span>
            </div>
          </div>

          <div className="flex-1 overflow-y-auto px-4 pt-3">
            <div className="mb-3 flex items-center justify-between px-2">
              <span className="text-[10px] font-semibold uppercase tracking-[0.17em] text-slate-500">
                Recentes
              </span>
              <MoreHorizontal size={17} className="text-slate-600" />
            </div>

            <div className="space-y-1">
              {conversations.map((conversation) => (
                <button
                  key={conversation.id}
                  onClick={() => selectConversation(conversation)}
                  className={`flex w-full items-center gap-3 rounded-lg px-3 py-3 text-left text-xs transition ${
                    activeConversation?.id === conversation.id
                      ? 'bg-cyan-400/10 text-cyan-300'
                      : 'text-slate-400 hover:bg-white/[0.04] hover:text-slate-200'
                  }`}
                >
                  <MessageSquare size={15} className="shrink-0" />
                  <span className="truncate">{conversation.title}</span>
                </button>
              ))}
            </div>

            <div className="mb-3 mt-9 px-2 text-[10px] font-semibold uppercase tracking-[0.17em] text-slate-500">
              Espaço de trabalho
            </div>

            <button
              onClick={() => setNotice('O gerenciamento de documentos será conectado ao backend.')}
              className="flex w-full items-center gap-3 rounded-lg px-3 py-3 text-xs text-slate-400 transition hover:bg-white/[0.04] hover:text-slate-200"
            >
              <FolderOpen size={16} />
              Documentos
            </button>
          </div>

          <div className="border-t border-white/[0.07] p-4">
            <div className="mb-3 flex items-center gap-3 rounded-xl bg-white/[0.035] p-3">
              <div className="flex h-9 w-9 items-center justify-center rounded-full bg-gradient-to-br from-cyan-400 to-blue-600 text-sm font-bold text-slate-950">
                N
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-xs font-medium text-slate-200">NetWatch Workspace</p>
                <p className="mt-1 text-[10px] text-slate-500">Assistente técnico</p>
              </div>
              <ChevronDown size={15} className="text-slate-500" />
            </div>

            <button
              onClick={() => setNotice('Configurações serão implementadas posteriormente.')}
              className="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-xs text-slate-500 transition hover:bg-white/[0.04] hover:text-slate-200"
            >
              <Settings size={15} />
              Configurações
            </button>
            <button
              onClick={() => setNotice('Central de ajuda em desenvolvimento.')}
              className="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-xs text-slate-500 transition hover:bg-white/[0.04] hover:text-slate-200"
            >
              <CircleHelp size={15} />
              Ajuda
            </button>
          </div>
        </aside>
      )}

      <main className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-[76px] shrink-0 items-center justify-between border-b border-white/[0.07] px-6 max-sm:px-4">
          <div className="flex items-center gap-3">
            {!sidebarOpen && (
              <button
                onClick={() => setSidebarOpen(true)}
                className="rounded-lg p-2 text-slate-400 transition hover:bg-white/5 hover:text-white"
                aria-label="Abrir menu"
              >
                <Menu size={19} />
              </button>
            )}
            <div>
              <h2 className="text-sm font-semibold text-slate-100">
                Assistente de documentação
              </h2>
              <p className="mt-1 text-[11px] text-slate-500">
                Consulte seus documentos técnicos com IA
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2 rounded-full border border-emerald-400/15 bg-emerald-400/[0.06] px-3 py-1.5">
            <span className="h-1.5 w-1.5 rounded-full bg-emerald-400" />
            <span className="text-[10px] font-medium text-emerald-300">
              Ambiente de desenvolvimento
            </span>
          </div>
        </header>

        <section className="flex flex-1 flex-col overflow-hidden">
          <div className="flex-1 overflow-y-auto">
            {messages.length === 0 ? (
              <div className="mx-auto flex min-h-full w-full max-w-4xl flex-col justify-center px-6 py-12 max-sm:px-4">
                <div className="mb-7 flex h-14 w-14 items-center justify-center rounded-2xl border border-cyan-400/20 bg-cyan-400/[0.08] text-cyan-300 shadow-[0_0_35px_rgba(34,211,238,0.08)]">
                  <Sparkles size={26} />
                </div>

                <p className="mb-3 text-xs font-medium uppercase tracking-[0.2em] text-cyan-400">
                  NETWATCH AI ASSISTANT
                </p>

                <h3 className="max-w-2xl text-4xl font-semibold leading-[1.18] tracking-tight text-white max-sm:text-3xl">
                  Inteligência para sua
                  <span className="text-cyan-400"> infraestrutura.</span>
                </h3>

                <p className="mt-4 max-w-xl text-sm leading-7 text-slate-400">
                  Consulte manuais, encontre especificações técnicas e obtenha
                  respostas contextualizadas a partir dos seus documentos.
                </p>

                <div className="mt-10 grid grid-cols-3 gap-3 max-md:grid-cols-1">
                  {suggestions.map((suggestion) => {
                    const Icon = suggestion.icon

                    return (
                      <button
                        key={suggestion.title}
                        onClick={() => setQuestion(suggestion.prompt)}
                        className="group rounded-xl border border-white/[0.08] bg-white/[0.025] p-4 text-left transition hover:border-cyan-400/30 hover:bg-cyan-400/[0.04]"
                      >
                        <div className="mb-4 flex h-9 w-9 items-center justify-center rounded-lg bg-cyan-400/10 text-cyan-300 transition group-hover:bg-cyan-400/15">
                          <Icon size={18} />
                        </div>
                        <p className="text-xs font-semibold text-slate-200">
                          {suggestion.title}
                        </p>
                        <p className="mt-2 text-[11px] leading-5 text-slate-500">
                          {suggestion.description}
                        </p>
                      </button>
                    )
                  })}
                </div>

                <div className="mt-9 flex items-center gap-2 text-[11px] text-slate-600">
                  <ShieldCheck size={14} />
                  Respostas baseadas nos documentos disponíveis
                </div>
              </div>
            ) : (
              <div className="mx-auto w-full max-w-3xl space-y-7 px-6 py-10 max-sm:px-4">
                {messages.map((message) => (
                  <div
                    key={message.id}
                    className={`flex gap-3 ${
                      message.role === 'user' ? 'justify-end' : 'justify-start'
                    }`}
                  >
                    {message.role === 'assistant' && (
                      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-cyan-400/10 text-cyan-300">
                        <Bot size={17} />
                      </div>
                    )}

                    <div
                      className={`max-w-[85%] rounded-2xl px-4 py-3 text-sm leading-7 ${
                        message.role === 'user'
                          ? 'bg-cyan-400 text-slate-950'
                          : 'border border-white/[0.08] bg-white/[0.035] text-slate-200'
                      }`}
                    >
                      {message.content}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="mx-auto w-full max-w-4xl shrink-0 px-6 pb-5 max-sm:px-4">
            {notice && (
              <div className="mb-3 flex items-start justify-between gap-3 rounded-lg border border-cyan-400/15 bg-cyan-400/[0.05] px-3 py-2.5 text-xs leading-5 text-cyan-200">
                <span>{notice}</span>
                <button
                  onClick={() => setNotice('')}
                  className="shrink-0 text-cyan-400 hover:text-white"
                  aria-label="Fechar aviso"
                >
                  <X size={15} />
                </button>
              </div>
            )}

            {selectedFile && (
              <div className="mb-3 flex items-center gap-2 rounded-lg border border-white/[0.08] bg-white/[0.03] px-3 py-2 text-xs text-slate-300">
                <FileText size={15} className="text-cyan-400" />
                <span className="flex-1 truncate">{selectedFile.name}</span>
                <button
                  onClick={() => setSelectedFile(null)}
                  className="text-slate-500 hover:text-white"
                  aria-label="Remover arquivo selecionado"
                >
                  <X size={15} />
                </button>
              </div>
            )}

            <form
              onSubmit={handleSubmit}
              className="rounded-2xl border border-white/[0.11] bg-[#111927] p-3 shadow-2xl shadow-black/20 transition focus-within:border-cyan-400/40"
            >
              <textarea
                value={question}
                onChange={(event) => setQuestion(event.target.value)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter' && !event.shiftKey) {
                    event.preventDefault()
                    event.currentTarget.form.requestSubmit()
                  }
                }}
                placeholder="Pergunte sobre seus equipamentos ou documentos..."
                rows={2}
                className="w-full resize-none bg-transparent px-2 py-2 text-sm leading-6 text-slate-100 outline-none placeholder:text-slate-600"
              />

              <div className="flex items-center justify-between pt-2">
                <div className="flex items-center gap-1">
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept=".pdf,application/pdf"
                    onChange={handleFileChange}
                    className="hidden"
                  />
                  <button
                    type="button"
                    onClick={() => fileInputRef.current?.click()}
                    className="flex items-center gap-2 rounded-lg px-3 py-2 text-xs text-slate-400 transition hover:bg-white/[0.06] hover:text-cyan-300"
                  >
                    <Paperclip size={15} />
                    Anexar PDF
                  </button>
                  <span className="hidden text-[10px] text-slate-600 sm:inline">
                    PDF até o limite configurado pela API
                  </span>
                </div>

                <button
                  type="submit"
                  disabled={!question.trim()}
                  className="flex h-9 w-9 items-center justify-center rounded-xl bg-cyan-400 text-slate-950 transition hover:bg-cyan-300 disabled:cursor-not-allowed disabled:opacity-30"
                  aria-label="Enviar mensagem"
                >
                  <Send size={16} />
                </button>
              </div>
            </form>

            <p className="mt-3 text-center text-[10px] text-slate-600">
              O NetWatch AI pode cometer erros. Verifique as informações técnicas nas fontes.
            </p>
          </div>
        </section>
      </main>
    </div>
  )
}

export default App
