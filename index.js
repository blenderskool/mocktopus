#!/usr/bin/env node

import { program } from 'commander';
import inquirer from 'inquirer';
import proto from 'proto-parser';
import chalk from 'chalk';
import path from 'path';
import ora from 'ora';
import fs from 'fs/promises';
import {
  askGPT,
  extractMessageDefinitions,
  getDefStr,
  sleep,
} from './utils.js';

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
  .option('-c, --code', 'generate code for generating mock data')
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

      let result;
      if (code) {
        result = await askGPT(
          `Generate JS code with "@faker-js/faker" library to create mock data for the "${definition}" proto definition in object format. Use only UUID for id fields and working image urls if needed\n\n${defStr}`
        );
      } else {
        result = await askGPT(
          `Generate valid JSON array with ${count} unique items and each item satisfying the "${definition}" proto definition. Use only UUID for id fields and working image urls if needed\n\n${defStr}`
        );
      }

      spinner.stop();
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

      const result = await askGPT(
        `Generate ${count} placeholder data for ${placeholder}`
      );

      spinner.stop();

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

      const result = await askGPT(
        `Generate tests code for the following code snippet based on what it does in the same language\n\n${inputStr}`
      );

      spinner.stop();

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

program
  .command('persona')
  .description('generate user personas for a product')
  .action(async () => {
    try {
      const { product } = await inquirer.prompt([
        {
          name: 'product',
          message: 'Describe your product:',
        },
      ]);

      const spinner = ora({
        text: 'Generating user personas for the product 🪄',
      }).start();

      const result = await askGPT(
        `Create a few user personas with name alliterations and different backgrounds for ${product}. Also add behavior, needs and wants, demographics to each persona`
      );

      spinner.stop();

      console.log();
      console.log(chalk.green('✅ User personas generated successfully 🐙'));
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

program.parse();
