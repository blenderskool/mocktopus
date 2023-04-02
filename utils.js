import { Configuration, OpenAIApi } from 'openai';

if (!process.env.MOCKTOPUS_OPENAI_KEY) {
  console.log(
    chalk.red(
      'Please add your OpenAI API key as env variable named "MOCKTOPUS_OPENAI_KEY"'
    )
  );
  process.exit(0);
}

const openai = new OpenAIApi(
  new Configuration({
    apiKey: process.env.MOCKTOPUS_OPENAI_KEY,
  })
);

// Extracts all nested "MessageDefinition" types from proto AST
export const extractMessageDefinitions = (def, result) => {
  if (!def) return;

  if (def?.syntaxType === 'MessageDefinition') {
    result.push(def);
  } else if (def.nested) {
    Object.values(def.nested).forEach((child) =>
      extractMessageDefinitions(child, result)
    );
  }
};

// Resolves dependencies and converts proto definition AST to string
export const getDefStr = (definitions, definition) => {
  const def = definitions.find(({ name }) => name === definition);
  const fields = Object.values(def.fields);

  const defStr = [];
  const fieldsStr = fields
    .map((field) => {
      // If the type is a user-defined type, include as a dependency
      if (field.type.syntaxType === 'Identifier') {
        defStr.push(getDefStr(definitions, field.type.value));
      }

      const baseStr = `${field.type.value} ${field.name}=${field.id};`;
      if (field.repeated) {
        return `repeated ${baseStr}`;
      }

      return baseStr;
    })
    .join('\n');

  defStr.push(`message ${def.name} {\n${fieldsStr}\n}`);

  return defStr.join('\n\n');
};

// Utility that resolves after a timeout
export const sleep = (duration) =>
  new Promise((resolve) => setTimeout(resolve, duration));

// Sends a message to ChatGPT and returns a response
export const askGPT = async (message) => {
  const response = await openai.createChatCompletion({
    model: 'gpt-3.5-turbo',
    messages: [
      {
        role: 'user',
        content: message,
      },
    ],
  });

  return response.choices[0].message.content;
};
