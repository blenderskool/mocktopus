#!/usr/bin/env node

import { program } from 'commander';
import inquirer from 'inquirer';
import proto from 'proto-parser';
import chalk from 'chalk';
import path from 'path';
import ora from 'ora';
import fs from 'fs/promises';
import { OpenAI } from 'node-openai';

if (!process.env.MOCKTOPUS_OPENAI_KEY) {
  console.log(
    chalk.red(
      'Please add your OpenAI API key as env variable named "MOCKTOPUS_OPENAI_KEY"'
    )
  );
  process.exit(0);
}

const openai = new OpenAI({
  apiKey: process.env.MOCKTOPUS_OPENAI_KEY,
}).v1();

// Extracts all nested "MessageDefinition" types from proto AST
const extractMessageDefinitions = (def, result) => {
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
const getDefStr = (definitions, definition) => {
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
const sleep = (duration) =>
  new Promise((resolve) => setTimeout(resolve, duration));

program
  .name('mocktopus')
  .description(
    `🐙 ${chalk
      .hex('#21D3A8')
      .bold('GPT powered')} CLI tool to generate mocks for anything!`
  )
  .version('1.0.0');

program
  .command('proto <source> <destination>')
  .option('-c --code', 'generate code for generating mock data')
  .description(
    'generate mock data for complex structures by analyzing proto definitions'
  )
  .action(async (inputPath, outputPath, { code }) => {
    console.log();

    const spinner = ora({ text: 'Scanning for definitions' }).start();
    const inputFile = (await fs.readFile(path.resolve(inputPath))).toString();
    await sleep(1000);

    spinner.stop();

    const protoDef = proto.parse(inputFile);

    const definitions = [];
    extractMessageDefinitions(protoDef.root, definitions);
    if (definitions.length === 0) {
      console.log(chalk.red('No definitions found, exiting...'));
      return;
    }

    console.log(chalk.green.bold(`${definitions.length} definitions found`));

    const { definition, count } = await inquirer.prompt([
      {
        name: 'definition',
        message: 'Which definition do you want mock data for?',
        type: 'list',
        choices: definitions.map(({ name, fields }) => ({
          value: name,
          name: `${name} (${Object.keys(fields).length} fields)`,
        })),
      },
      // The following question is only relevant when code is not being generated
      {
        name: 'count',
        message: 'Number of records to generate?',
        type: 'number',
        default: 1,
        when: () => !code,
      },
    ]);

    try {
      const defStr = getDefStr(definitions, definition);
      const spinner = ora({
        text: code
          ? 'Generating code for generating mock data for proto definition 🪄'
          : 'Generating mock data for proto definition 🪄',
      }).start();

      let response;
      if (code) {
        response = await openai.chat.create({
          model: 'gpt-3.5-turbo',
          messages: [
            {
              role: 'user',
              content: `Generate JS code with "@faker-js/faker" library to create mock data for the "${definition}" proto definition in object format. Use only UUID for id fields if needed\n\n${defStr}`,
            },
          ],
        });
      } else {
        response = await openai.chat.create({
          model: 'gpt-3.5-turbo',
          messages: [
            {
              role: 'user',
              content: `Generate ${count} unique array items with mock data in JSON format for the "${definition}" proto definition. Use only UUID for id fields if needed\n\n${defStr}`,
            },
          ],
        });
      }

      spinner.stop();
      const result = response.choices[0].message.content;
      await fs.writeFile(path.resolve(outputPath), result);

      console.log();
      console.log(
        chalk.green(
          code
            ? '✅ Code for mock data generated successfully 🐙'
            : '✅ Mock data generated successfully 🐙'
        )
      );
    } catch (err) {
      console.log();
      console.log(
        chalk.red.bold('⚠️ Unexpected error occurred, try different definition')
      );
      console.log(chalk.white.dim(err));
    }
  });

program
  .command('placeholder')
  .description('generate mock data from natural descriptions')
  .action(async () => {
    const { placeholder, count } = await inquirer.prompt([
      {
        name: 'placeholder',
        message: 'What do you want placeholder for?',
      },
      {
        name: 'count',
        message: 'Number of records to generate?',
        type: 'number',
        default: 1,
      },
    ]);

    try {
      const spinner = ora({
        text: 'Generating mock placeholder data 🪄',
      }).start();

      const response = await openai.chat.create({
        model: 'gpt-3.5-turbo',
        messages: [
          {
            role: 'user',
            content: `Generate ${count} placeholder data for ${placeholder}`,
          },
        ],
      });

      spinner.stop();

      const result = response.choices[0].message.content;
      console.log();
      console.log(chalk.green('✅ Mock data generated successfully 🐙'));
      console.log();
      console.log(result);
    } catch (err) {
      console.log();
      console.log(
        chalk.red.bold('⚠️ Unexpected error occurred, try again later')
      );
      console.log(chalk.white.dim(err));
    }
  });

program
  .command('tests <source> <destination>')
  .description('generate test cases for code snippets')
  .action(async (source, outputPath) => {
    try {
      let inputPath = source;
      let range = [0, Number.MAX_SAFE_INTEGER];
      if (source.includes('#')) {
        [inputPath, range] = inputPath.split('#');
        range = range.split(':');
      }

      const inputFile = (await fs.readFile(path.resolve(inputPath))).toString();
      const inputStr = inputFile
        .split('\n')
        .slice(range[0] - 1, range[1])
        .join('\n');

      const spinner = ora({
        text: 'Generating tests for code snippet 🪄',
      }).start();

      const response = await openai.chat.create({
        model: 'gpt-3.5-turbo',
        messages: [
          {
            role: 'user',
            content: `Generate tests code for the following code snippet based on what it does in the same language\n\n${inputStr}`,
          },
        ],
      });

      spinner.stop();

      const result = response.choices[0].message.content;
      await fs.writeFile(path.resolve(outputPath), result);

      console.log();
      console.log(chalk.green('✅ Test cases generated successfully 🐙'));
    } catch (err) {
      console.log();
      console.log(
        chalk.red.bold(
          '⚠️ Unexpected error occurred, try with different code snippet'
        )
      );
      console.log(chalk.white.dim(err));
    }
  });

program.parse();
